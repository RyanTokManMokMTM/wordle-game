package gameClient

import "net"

type IGameClient interface {
	HandleRequest()
	SetWordHistory(string)
	SetGuessingWord()

	GetTotalRound() uint
	GetWordList() []string
	GetGuessingWord() string
	GetWordHistory() []string
	GetConn() net.Conn

	Reset()
}
