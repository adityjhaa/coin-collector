package client

import (
	"coin-collector/common"
	"encoding/binary"
	"log"
	"math"
	"math/rand"
	"net"
)

type Network struct {
	conn     *net.UDPConn
	server   *net.UDPAddr
	PlayerID common.PlayerID

	StateChan chan []byte // raw server world state packets
}

func NewNetwork() *Network {
	return &Network{
		StateChan: make(chan []byte, 32),
	}
}

func (n *Network) Connect() error {
	addr, err := net.ResolveUDPAddr("udp", common.ServerPort)
	if err != nil {
		return err
	}
	n.server = addr

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	n.conn = conn

	// Send CONNECT
	n.conn.Write([]byte{common.MsgConnect})

	// Wait for WELCOME
	buf := make([]byte, 8)
	for {
		nBytes, _, err := n.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if nBytes >= 2 && buf[0] == common.MsgConnect {
			n.PlayerID = common.PlayerID(buf[1])
			break
		}
	}

	log.Println("Connected. PlayerID =", n.PlayerID)

	// Send SPAWN
	x := rand.Float32() * 800
	y := rand.Float32() * 600

	spawn := make([]byte, 1+4+4)
	spawn[0] = common.MsgSpawn
	binary.LittleEndian.PutUint32(spawn[1:5], math.Float32bits(x))
	binary.LittleEndian.PutUint32(spawn[5:9], math.Float32bits(y))

	n.conn.Write(spawn)

	// Start async listener
	go n.listenLoop()

	return nil
}

func (n *Network) listenLoop() {
	buf := make([]byte, 2048)

	for {
		nBytes, _, err := n.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		if nBytes == 0 {
			continue
		}

		msgType := buf[0]

		if msgType == common.MsgState {
			// push raw world state packet
			data := make([]byte, nBytes)
			copy(data, buf[:nBytes])
			n.StateChan <- data
		}
	}
}

func (n *Network) SendInput(mask uint8) {
	pkt := make([]byte, 1+2+1) // type + playerID + mask
	pkt[0] = common.MsgInput
	binary.LittleEndian.PutUint16(pkt[1:3], uint16(n.PlayerID))
	pkt[3] = mask

	n.conn.Write(pkt)

}
