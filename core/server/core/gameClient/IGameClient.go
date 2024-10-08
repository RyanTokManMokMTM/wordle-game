package gameClient

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"net"
)

type IGameClient interface {
	Run()
	read()
	write()
	Closed()

	SetUserName(name string)
	GetConn() net.Conn
	GetClientId() string
	GetName() string

	GetClosedEvent() chan struct{}
	GetMessage() chan packet.BasicPacket
	GetGameGuessingInput() chan []byte

	SendToClient(code uint, message, pkType string, data []byte)
}
