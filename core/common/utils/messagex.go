package utils

import (
	"encoding/binary"
	"net"
)

func SendMessage(conn net.Conn, pkData []byte) error {

	msgLen := uint32(len(pkData))
	err := binary.Write(conn, binary.BigEndian, msgLen)
	if err != nil {
		return err
	}

	_, err = conn.Write(pkData)
	if err != nil {
		return err
	}

	return nil
}
