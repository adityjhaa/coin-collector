package client

import (
	"encoding/binary"
	"math"
	"time"

	"coin-collector/common"
)

type PlayerSnapshot struct {
	ID    common.PlayerID
	X, Y  float32
	Score uint16
}

type WorldSnapshot struct {
	Players   []PlayerSnapshot
	Coins     []common.Coin
	Timestamp int64
}

func NowMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func ParseWorldState(data []byte) WorldSnapshot {
	offset := 1

	if len(data) < offset+2 {
		return WorldSnapshot{Timestamp: NowMs()}
	}
	numPlayers := int(binary.LittleEndian.Uint16(data[offset:]))
	offset += 2

	players := make([]PlayerSnapshot, 0, numPlayers)
	for i := 0; i < numPlayers; i++ {
		if len(data) < offset+2+4+4+2 {
			break
		}
		id := common.PlayerID(binary.LittleEndian.Uint16(data[offset:]))
		offset += 2
		xbits := binary.LittleEndian.Uint32(data[offset:])
		offset += 4
		ybits := binary.LittleEndian.Uint32(data[offset:])
		offset += 4
		score := binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		players = append(players, PlayerSnapshot{
			ID:    id,
			X:     math.Float32frombits(xbits),
			Y:     math.Float32frombits(ybits),
			Score: score,
		})
	}

	if len(data) < offset+2 {
		return WorldSnapshot{Players: players, Coins: nil, Timestamp: NowMs()}
	}
	numCoins := int(binary.LittleEndian.Uint16(data[offset:]))
	offset += 2

	coins := make([]common.Coin, 0, numCoins)
	for i := 0; i < numCoins; i++ {
		if len(data) < offset+2+4+4 {
			break
		}
		cid := int(binary.LittleEndian.Uint16(data[offset:]))
		offset += 2
		xbits := binary.LittleEndian.Uint32(data[offset:])
		offset += 4
		ybits := binary.LittleEndian.Uint32(data[offset:])
		offset += 4
		coins = append(coins, common.Coin{
			ID: cid,
			X:  math.Float32frombits(xbits),
			Y:  math.Float32frombits(ybits),
		})
	}

	return WorldSnapshot{
		Players:   players,
		Coins:     coins,
		Timestamp: NowMs(),
	}
}
