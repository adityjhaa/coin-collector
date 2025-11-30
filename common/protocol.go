package common

const (
	MsgConnect = iota
	MsgInput
	MsgState
	MsgSpawn
)

type InputMessage struct {
	ID    PlayerID
	Input uint8
}

type SpawnMessage struct {
	X float32
	Y float32
}
