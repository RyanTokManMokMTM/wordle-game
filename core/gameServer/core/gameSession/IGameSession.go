package gameSession

import "github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"

type IGameSession interface {
	SetJoinedPlayer(uid string, player gameClient.IGameClient)
	SetSessionStatus(status string)

	GetSessionId() string
	GetAllJoinedPlayer() []gameClient.IGameClient
	GetSessionStatus() string
}
