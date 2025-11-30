package common

import "net"

type PlayerID uint16

type Player struct {
	ID        PlayerID
	X, Y      float32
	Score     uint16
	Addr      *net.UDPAddr
	LastInput uint8
	Spawned   bool
	LastHeard int64
}

type Coin struct {
	ID int
	X  float32
	Y  float32
}
