package server

import (
	"bytes"
	"coin-collector/common"
	"encoding/binary"
)

func (s *Server) BuildWorldState() []byte {
	buf := new(bytes.Buffer)

	// Header: message type
	binary.Write(buf, binary.LittleEndian, uint8(common.MsgState))

	// Number of players
	binary.Write(buf, binary.LittleEndian, uint16(len(s.players)))

	for _, p := range s.players {
		binary.Write(buf, binary.LittleEndian, uint16(p.ID))
		binary.Write(buf, binary.LittleEndian, float32(p.X))
		binary.Write(buf, binary.LittleEndian, float32(p.Y))
		binary.Write(buf, binary.LittleEndian, uint16(p.Score))
	}

	// Number of coins
	binary.Write(buf, binary.LittleEndian, uint16(len(s.coins)))

	for _, c := range s.coins {
		binary.Write(buf, binary.LittleEndian, uint16(c.ID))
		binary.Write(buf, binary.LittleEndian, float32(c.X))
		binary.Write(buf, binary.LittleEndian, float32(c.Y))
	}

	return buf.Bytes()
}
