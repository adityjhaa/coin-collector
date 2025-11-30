package server

import (
	"log"
	"sync"
	"time"

	"coin-collector/common"
)

type Server struct {
	conn *netUDPWrapper

	lock sync.RWMutex

	players map[common.PlayerID]*common.Player
	coins   []common.Coin

	lastCoinSpawn time.Time
	nextPlayerID  common.PlayerID
	nextCoinID    int
}

func NewServer() *Server {
	return &Server{
		players:       make(map[common.PlayerID]*common.Player),
		coins:         make([]common.Coin, 0),
		lastCoinSpawn: time.Now(),
		nextPlayerID:  1,
		nextCoinID:    1,
	}
}

func (s *Server) StartGameLoop() {
	ticker := time.NewTicker(time.Duration(common.TickDelayMs) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.applyPlayerInputs()
		s.resolvePlayerCollisions()
		s.spawnCoin()
		s.checkCoinPickups()
		s.removeTimedOut()
		s.clearCoins()

		state := s.buildWorldState()
		s.broadcastState(state)
	}
}

func (s *Server) removeTimedOut() {
	now := NowMs()

	s.lock.Lock()
	defer s.lock.Unlock()

	for id, p := range s.players {
		if now-p.LastHeard > 3000 {
			log.Println("Removing timed-out player:", id)
			delete(s.players, id)
		}
	}
}

func (s *Server) clearCoins() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.players) == 0 {
		s.coins = nil
		s.nextCoinID = 1
		s.nextPlayerID = 1
	}
}

func NowMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
