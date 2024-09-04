package manager

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"
)

type IGameClientManager interface {
	SetGameClient(clientId string, client gameClient.IGameClient)
	RemoveGameClient(clientId string)

	GetGameClient(clientId string) (gameClient.IGameClient, bool)
	GetGameClientList() []gameClient.IGameClient
}
