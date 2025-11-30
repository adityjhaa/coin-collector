package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"coin-collector/common"
)

type Network struct {
	conn     *net.UDPConn
	server   *net.UDPAddr
	PlayerID common.PlayerID

	StateChan chan []byte
	closed    chan struct{}
}

func NewNetwork() *Network {
	return &Network{
		StateChan: make(chan []byte, 64),
		closed:    make(chan struct{}),
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
	n.conn.SetReadDeadline(time.Time{})

	if _, err := n.conn.Write([]byte{byte(common.MsgConnect)}); err != nil {
		return err
	}

	buf := make([]byte, 64)
	n.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	nBytes, _, err := n.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	n.conn.SetReadDeadline(time.Time{})

	if nBytes >= 3 && buf[0] == byte(common.MsgWelcome) {
		id := common.PlayerID(binary.LittleEndian.Uint16(buf[1:3]))
		n.PlayerID = id
		log.Println("Received WELCOME id=", id)
	} else {
		return fmt.Errorf("invalid WELCOME packet")
	}

	go n.listenLoop()
	return nil
}

func (n *Network) listenLoop() {
	buf := make([]byte, 4096)
	for {
		nBytes, _, err := n.conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-n.closed:
				return
			default:
			}
			continue
		}
		if nBytes == 0 {
			continue
		}
		if buf[0] == byte(common.MsgState) {
			packet := make([]byte, nBytes)
			copy(packet, buf[:nBytes])
			select {
			case n.StateChan <- packet:
			default:
				select {
				case <-n.StateChan:
				default:
				}
				select {
				case n.StateChan <- packet:
				default:
				}
			}
		}
	}
}

func (n *Network) SendInput(mask uint8) {
	pkt := make([]byte, 1+2+1)
	pkt[0] = byte(common.MsgInput)
	binary.LittleEndian.PutUint16(pkt[1:3], uint16(n.PlayerID))
	pkt[3] = mask
	_, _ = n.conn.Write(pkt)
}

func (n *Network) Close() {
	close(n.closed)
	if n.conn != nil {
		_ = n.conn.Close()
	}
}
