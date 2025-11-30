package server

import (
	"coin-collector/common"
	"encoding/binary"
	"log"
	"math"
	"net"
	"time"
)

type Server struct {
	conn          *net.UDPConn
	players       map[common.PlayerID]*common.Player
	coins         []common.Coin
	lastCoinSpawn time.Time
}

func NewServer() *Server {
	return &Server{
		players:       make(map[common.PlayerID]*common.Player),
		coins:         make([]common.Coin, 0),
		lastCoinSpawn: time.Now(),
	}
}

func (s *Server) StartNetwork() error {
	addr, err := net.ResolveUDPAddr("udp", common.ServerPort)
	if err != nil {
		return err
	}

	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	log.Println("Server listening on", common.ServerPort)

	go s.listenLoop()

	return nil
}

func (s *Server) listenLoop() {
	buf := make([]byte, 1024)

	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Read error:", err)
			continue
		}

		time.Sleep(200 * time.Millisecond)

		if n == 0 {
			continue
		}
		msgType := buf[0]

		switch msgType {

		case common.MsgConnect:
			id := s.AddPlayer(addr)
			s.sendWelcome(addr, id)
			log.Println("Player connected:", id, addr)

		case common.MsgSpawn:
			s.handleSpawn(addr, buf[:n])

		case common.MsgInput:
			s.handleInput(addr, buf[:n])
		}
	}
}

func (s *Server) sendWelcome(addr *net.UDPAddr, id common.PlayerID) {
	data := []byte{
		common.MsgConnect,
		byte(id),
	}
	s.conn.WriteToUDP(data, addr)
}

func (s *Server) handleSpawn(addr *net.UDPAddr, packet []byte) {
	pid := s.playerIDFromAddr(addr)
	if pid == 0 {
		return
	}

	x := math.Float32frombits(binary.LittleEndian.Uint32(packet[1:5]))
	y := math.Float32frombits(binary.LittleEndian.Uint32(packet[5:9]))

	s.SetPlayerSpawn(pid, x, y)
	log.Println("Player", pid, "spawned at", x, y)
}

func (s *Server) handleInput(addr *net.UDPAddr, packet []byte) {
    if len(packet) < 4 {
        return
    }

    pid := common.PlayerID(binary.LittleEndian.Uint16(packet[1:3]))
    mask := packet[3]

    p, ok := s.players[pid]
    if !ok { return }

    p.LastInput = mask
}


func (s *Server) Broadcast(data []byte) {
	for _, p := range s.players {
		time.Sleep(200 * time.Millisecond)
		s.conn.WriteToUDP(data, p.Addr)
	}
}

func (s *Server) playerIDFromAddr(addr *net.UDPAddr) common.PlayerID {
	for id, p := range s.players {
		if p.Addr.String() == addr.String() {
			return id
		}
	}
	return 0
}
