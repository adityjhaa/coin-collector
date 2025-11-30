package common

const (
	ServerPort = ":8080"

	TickRate    = 30
	TickDelayMs = 1000 / TickRate

	MoveSpeed    = 200.0 / TickRate
	PickupRadius = 24.0

	CoinSpawnIntervalSec = 7
	MaxCoins             = 10

	SimulatedLatencyMs = 200
)
