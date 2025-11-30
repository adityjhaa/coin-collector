package server

import (
	"bytes"
	"encoding/binary"

	"coin-collector/common"
)

func (s *Server) buildWorldState() []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()

	buf := new(bytes.Buffer)
	buf.WriteByte(byte(common.MsgState))

	binary.Write(buf, binary.LittleEndian, uint16(len(s.players)))

	for _, p := range s.players {
		binary.Write(buf, binary.LittleEndian, uint16(p.ID))
		binary.Write(buf, binary.LittleEndian, p.X)
		binary.Write(buf, binary.LittleEndian, p.Y)
		binary.Write(buf, binary.LittleEndian, p.Score)
	}

	binary.Write(buf, binary.LittleEndian, uint16(len(s.coins)))
	for _, c := range s.coins {
		binary.Write(buf, binary.LittleEndian, uint16(c.ID))
		binary.Write(buf, binary.LittleEndian, c.X)
		binary.Write(buf, binary.LittleEndian, c.Y)
	}

	return buf.Bytes()
}
