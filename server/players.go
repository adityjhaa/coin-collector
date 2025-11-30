package server

import (
	"coin-collector/common"
	"net"
	"sync/atomic"
)

var nextPlayerID uint32 = 1

func (s *Server) AddPlayer(addr *net.UDPAddr) common.PlayerID {
	id := common.PlayerID(atomic.AddUint32(&nextPlayerID, 1))

	s.players[id] = &common.Player{
		ID:        id,
		X:         0,
		Y:         0,
		Score:     0,
		Addr:      addr,
		LastInput: 0,
	}

	return id
}

func (s *Server) SetPlayerSpawn(id common.PlayerID, x, y float32) {
	if p, ok := s.players[id]; ok {
		p.X = float64(x)
		p.Y = float64(y)
	}
}

func (s *Server) RemovePlayer(id common.PlayerID) {
	delete(s.players, id)
}

func (s *Server) ApplyPlayerInputs() {
	for _, p := range s.players {
		if p.LastInput == 0 {
			continue
		}

		// bitmask: W=1, A=2, S=4, D=8
		if p.LastInput&1 != 0 { // W
			p.Y -= common.MoveSpeed
		}
		if p.LastInput&2 != 0 { // A
			p.X -= common.MoveSpeed
		}
		if p.LastInput&4 != 0 { // S
			p.Y += common.MoveSpeed
		}
		if p.LastInput&8 != 0 { // D
			p.X += common.MoveSpeed
		}

		// clamp to screen for safety
		if p.X < 0 {
			p.X = 0
		}
		if p.Y < 0 {
			p.Y = 0
		}
		if p.X > 800 {
			p.X = 800
		}
		if p.Y > 600 {
			p.Y = 600
		}
	}
}
