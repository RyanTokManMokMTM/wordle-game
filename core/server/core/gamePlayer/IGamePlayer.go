package gamePlayer

import "github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameClient"

type IGamePlayer interface {
	GetClient() gameClient.IGameClient

	SetScore(score uint)
	GetScore() uint
}
