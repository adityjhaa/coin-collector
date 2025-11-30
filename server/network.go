package server

import (
	"encoding/binary"
	"log"
	"net"
	"time"

	"coin-collector/common"
)

type netUDPWrapper struct {
	conn *net.UDPConn
}

func (w *netUDPWrapper) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return w.conn.WriteToUDP(b, addr)
}

func (w *netUDPWrapper) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return w.conn.ReadFromUDP(b)
}

func (s *Server) StartNetwork() error {
	addr, err := net.ResolveUDPAddr("udp", common.ServerPort)
	if err != nil {
		return err
	}
	rawConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.conn = &netUDPWrapper{conn: rawConn}

	go s.listenLoop()
	return nil
}

func (s *Server) listenLoop() {
	buf := make([]byte, 2048)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("read err:", err)
			continue
		}
		if n == 0 {
			continue
		}

		msgType := buf[0]
		switch msgType {
		case common.MsgConnect:
			id := s.AddPlayer(addr)
			b := make([]byte, 1+2)
			b[0] = byte(common.MsgWelcome)
			binary.LittleEndian.PutUint16(b[1:3], uint16(id))
			go s.delayedSend(b, addr)
			log.Println("CONNECT from", addr, "-> id", id)

		case common.MsgInput:
			if n < 1+2+1 {
				continue
			}
			id := common.PlayerID(binary.LittleEndian.Uint16(buf[1:3]))
			mask := buf[3]
			s.lock.Lock()
			if p, ok := s.players[id]; ok {
				p.LastHeard = NowMs()
				p.LastInput = mask
			}
			s.lock.Unlock()

		default:
		}
	}
}

func (s *Server) delayedSend(data []byte, addr *net.UDPAddr) {
	lat := time.Duration(common.SimulatedLatencyMs) * time.Millisecond
	time.Sleep(lat)
	_, err := s.conn.WriteToUDP(data, addr)
	if err != nil {
		log.Println("delayed send err:", err)
	}
}

func (s *Server) broadcastState(state []byte) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, p := range s.players {
		if !p.Spawned {
			continue
		}
		packet := make([]byte, len(state))
		copy(packet, state)
		addr := p.Addr
		go s.delayedSend(packet, addr)
	}
}
