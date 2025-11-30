package server

import (
	"coin-collector/common"
	"math"
	"math/rand"
	"net"
)

const playerHalfSize = 13

func (s *Server) AddPlayer(addr *net.UDPAddr) common.PlayerID {
	s.lock.Lock()
	defer s.lock.Unlock()

	id := s.nextPlayerID
	s.nextPlayerID++

	x := float32(rand.Intn(760) + 20)
	y := float32(rand.Intn(560) + 20)

	p := &common.Player{
		ID:        id,
		X:         x,
		Y:         y,
		Score:     0,
		Addr:      addr,
		LastInput: 0,
		LastHeard: NowMs(),
		Spawned:   true,
	}
	s.players[id] = p
	return id
}

func (s *Server) applyPlayerInputs() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, p := range s.players {
		if !p.Spawned {
			continue
		}
		m := p.LastInput
		if m == 0 {
			continue
		}

		dx, dy := 0.0, 0.0

		if m&1 != 0 { // W
			dy -= 1
		}
		if m&2 != 0 { // A
			dx -= 1
		}
		if m&4 != 0 { // S
			dy += 1
		}
		if m&8 != 0 { // D
			dx += 1
		}

		if dx != 0 || dy != 0 {
			length := math.Hypot(dx, dy)
			dx /= length
			dy /= length

			p.X += float32(dx * common.MoveSpeed)
			p.Y += float32(dy * common.MoveSpeed)
		}

		if p.X < playerHalfSize {
			p.X = playerHalfSize
		}
		if p.Y < playerHalfSize {
			p.Y = playerHalfSize
		}
		if p.X > 800-playerHalfSize {
			p.X = 800 - playerHalfSize
		}
		if p.Y > 600-playerHalfSize {
			p.Y = 600 - playerHalfSize
		}
	}
}

func (s *Server) resolvePlayerCollisions() {
	s.lock.Lock()
	defer s.lock.Unlock()

	const minDist = 26.0

	for _, p1 := range s.players {
		for _, p2 := range s.players {
			if p1.ID >= p2.ID {
				continue
			}

			dx := float64(p2.X - p1.X)
			dy := float64(p2.Y - p1.Y)
			dist := math.Hypot(dx, dy)

			if dist < minDist && dist > 0 {
				overlap := minDist - dist
				nx := dx / dist
				ny := dy / dist

				p1.X -= float32(nx * overlap / 2)
				p1.Y -= float32(ny * overlap / 2)
				p2.X += float32(nx * overlap / 2)
				p2.Y += float32(ny * overlap / 2)
			}
		}
	}
}
