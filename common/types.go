package common

import "net"

type PlayerID int

type Player struct {
	ID        PlayerID
	X, Y      float64
	Score     int
	Addr      *net.UDPAddr
	LastInput uint8
}

type Coin struct {
	ID int
	X  float64
	Y  float64
}

type WorldState struct {
	Players []PlayerState
	Coins   []Coin
}

type PlayerState struct {
	ID    PlayerID
	X, Y  float64
	Score int
}
