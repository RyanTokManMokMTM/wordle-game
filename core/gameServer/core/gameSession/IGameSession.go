package gameSession

import "github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"

type IGameSession interface {
	SetJoinedPlayer(uid string, player gameClient.IGameClient)
	SetExitedPlayer(uid string)
	SetSessionStatus(status string)

	GetSessionId() string
	GetAllJoinedPlayer() []gameClient.IGameClient
	GetSessionStatus() string
	GetSessionHost() gameClient.IGameClient
	GetSessionName() string
	GetSessionPlayerInfo() (min uint, max uint, current uint)
	GetSessionWordList() []string
}
