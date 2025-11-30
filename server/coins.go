package server

import (
	"coin-collector/common"
	"math"
	"math/rand"
	"time"
)

var nextCoinID = 1

func (s *Server) MaybeSpawnCoin() {
	if time.Since(s.lastCoinSpawn).Seconds() >= common.CoinSpawnInterval {
		s.lastCoinSpawn = time.Now()

		coin := common.Coin{
			ID: nextCoinID,
			X:  float64(rand.Intn(780) + 20),
			Y:  float64(rand.Intn(580) + 20),
		}
		nextCoinID++

		s.coins = append(s.coins, coin)
	}
}

func (s *Server) CheckCoinPickups() {
	for pid, p := range s.players {
		remaining := s.coins[:0]

		for _, c := range s.coins {
			dx := p.X - c.X
			dy := p.Y - c.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < common.PickupRadius {
				p.Score++
				continue
			}

			remaining = append(remaining, c)
		}

		s.coins = remaining
		s.players[pid] = p
	}
}
