package gameClient

import "net"

type IGameClient interface {
	Read()
	Write()
	Closed()

	SetUserName(name string)
	GetConn() net.Conn
	GetClientId() string

	GetClosedEvent() chan struct{}
}
