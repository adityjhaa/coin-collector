package server

import (
	"coin-collector/common"
	"log"
	"time"
)

func (s *Server) StartGameLoop() {
	ticker := time.NewTicker(time.Duration(common.TickDelayMs) * time.Millisecond)
	defer ticker.Stop()

	log.Println("Server loop started")

	for range ticker.C {

		s.ApplyPlayerInputs()

		s.MaybeSpawnCoin()

		s.CheckCoinPickups()

		state := s.BuildWorldState()

		s.Broadcast(state)
	}
}
