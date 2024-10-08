package gameRoom

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gamePlayer"
)

type IGameRoom interface {
	AddPlayer(uid string, p gamePlayer.IGamePlayer)
	RemovePlayer(uid string) bool

	SetRoomStatus(status string)
	SetGuessingWord()
	StartGame(player gamePlayer.IGamePlayer)

	GetRoomId() string
	GetAllPlayer() []gamePlayer.IGamePlayer
	GetRoomStatus() string
	GetRoomHost() gamePlayer.IGamePlayer
	GetRoomName() string
	GetRoomWordList() []string

	RemoveAllPlayer()
	NotifyPlayerWithMessage(p gamePlayer.IGamePlayer, message string)
	GetTheGameIsOver() chan struct{}

	Close()
}
