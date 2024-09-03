package gameSession

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"
	"github.com/google/uuid"
	"sync"
)

type GameSession struct {
	sync.Mutex
	sessionId string
	host      gameClient.IGameClient
	name      string
	minPlayer uint
	maxPlayer uint

	status  string
	players map[string]gameClient.IGameClient
}

func NewGameSession(host gameClient.IGameClient, name string, minPlayer, maxPlayer uint) IGameSession {
	return &GameSession{
		sessionId: uuid.NewString(),
		host:      host,
		name:      name,
		minPlayer: minPlayer,
		maxPlayer: maxPlayer,
		players:   make(map[string]gameClient.IGameClient),
	}
}

func (gs *GameSession) SetJoinedPlayer(uid string, player gameClient.IGameClient) {
	gs.Lock()
	defer gs.Unlock()

	gs.players[uid] = player

}

func (gs *GameSession) SetSessionStatus(status string) {
	gs.status = status
}

func (gs *GameSession) GetAllJoinedPlayer() []gameClient.IGameClient {
	allPlayer := make([]gameClient.IGameClient, 0)
	for _, p := range gs.players {
		allPlayer = append(allPlayer, p)
	}
	return allPlayer
}

func (gs *GameSession) GetSessionStatus() string {
	return gs.status
}

func (gs *GameSession) GetSessionId() string {
	return gs.sessionId
}
