package common

const (
	ServerPort = ":8080"

	TickRate    = 30
	TickDelayMs = 1000 / TickRate

	MoveSpeed    = 200.0 / TickRate
	PickupRadius = 25.0

	CoinSpawnIntervalSec = 8
	MaxCoins             = 10

	SimulatedLatencyMs = 200
)
