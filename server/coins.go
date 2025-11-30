package server

import (
	"math"
	"math/rand"
	"time"

	"coin-collector/common"
)

func (s *Server) spawnCoin() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.players) == 0 {
		return
	}

	if len(s.coins) >= common.MaxCoins {
		return
	}

	if time.Since(s.lastCoinSpawn).Seconds() < common.CoinSpawnIntervalSec {
		return
	}

	s.lastCoinSpawn = time.Now()

	c := common.Coin{
		ID: s.nextCoinID,
		X:  float32(rand.Intn(760) + 20),
		Y:  float32(rand.Intn(560) + 20),
	}
	s.nextCoinID++
	s.coins = append(s.coins, c)
}

func (s *Server) checkCoinPickups() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.coins) == 0 {
		return
	}

	remaining := make([]common.Coin, 0, len(s.coins))
	for _, c := range s.coins {
		picked := false
		for _, p := range s.players {
			if !p.Spawned {
				continue
			}
			dx := float64(p.X) - float64(c.X)
			dy := float64(p.Y) - float64(c.Y)
			if math.Hypot(dx, dy) < common.PickupRadius {
				p.Score++
				picked = true
				break
			}
		}
		if !picked {
			remaining = append(remaining, c)
		}
	}
	s.coins = remaining
}
