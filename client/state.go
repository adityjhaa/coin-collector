package client

import (
	"encoding/binary"
	"coin-collector/common"
	"time"
)

type PlayerSnapshot struct {
	ID    common.PlayerID
	X, Y  float32
	Score int
}

type WorldSnapshot struct {
	Players   []PlayerSnapshot
	Coins     []common.Coin
	Timestamp int64 // capture time for interpolation
}

func ParseWorldState(data []byte) WorldSnapshot {
	offset := 1 // skip msgType

	// --- players ---
	numPlayers := int(binary.LittleEndian.Uint16(data[offset:]))
	offset += 2

	players := make([]PlayerSnapshot, 0, numPlayers)

	for i := 0; i < numPlayers; i++ {
		id := common.PlayerID(binary.LittleEndian.Uint16(data[offset:]))
		offset += 2

		x := binary.LittleEndian.Uint32(data[offset:])
		offset += 4

		y := binary.LittleEndian.Uint32(data[offset:])
		offset += 4

		score := binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		players = append(players, PlayerSnapshot{
			ID:    id,
			X:     float32(x),
			Y:     float32(y),
			Score: int(score),
		})
	}

	// --- coins ---
	numCoins := int(binary.LittleEndian.Uint16(data[offset:]))
	offset += 2

	coins := make([]common.Coin, 0, numCoins)

	for i := 0; i < numCoins; i++ {
		cid := int(binary.LittleEndian.Uint16(data[offset:]))
		offset += 2

		x := binary.LittleEndian.Uint32(data[offset:])
		offset += 4

		y := binary.LittleEndian.Uint32(data[offset:])
		offset += 4

		coins = append(coins, common.Coin{
			ID: cid,
			X:  float64(float32(x)),
			Y:  float64(float32(y)),
		})
	}

	return WorldSnapshot{
		Players:   players,
		Coins:     coins,
		Timestamp: NowMs(),
	}
}

func NowMs() int64 {
	return int64(float64(time.Now().UnixNano()) / 1e6)
}
