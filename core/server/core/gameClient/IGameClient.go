package gameClient

import "net"

type IGameClient interface {
	// HandleRequest handing request from client
	HandleRequest()
	SetGuessingWord()

	GetTotalRound() uint
	GetWordList() []string
	GetGuessingWord() string

	// GetConn get client connection
	GetConn() net.Conn
}
