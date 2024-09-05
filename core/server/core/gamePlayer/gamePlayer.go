package gamePlayer

import "github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameClient"

type GamePlayer struct {
	client gameClient.IGameClient
	score  uint
}

func NewPlayer(client gameClient.IGameClient) IGamePlayer {
	return &GamePlayer{
		client: client,
		score:  0,
	}
}

func (gp *GamePlayer) GetClient() gameClient.IGameClient {
	return gp.client
}

func (gp *GamePlayer) SetScore(score uint) {
	gp.score = score
}

func (gp *GamePlayer) GetScore() uint {
	return gp.score
}
