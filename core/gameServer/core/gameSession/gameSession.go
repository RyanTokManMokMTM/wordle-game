package gameSession

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/status"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"
	"github.com/google/uuid"
	"sync"
)

type GameSession struct {
	sync.Mutex
	sessionId     string
	host          gameClient.IGameClient
	name          string
	minPlayer     uint
	currentPlayer uint
	maxPlayer     uint
	wordList      []string
	status        string
	players       map[string]gameClient.IGameClient
}

func NewGameSession(host gameClient.IGameClient, name string, minPlayer, maxPlayer uint, wordList []string) IGameSession {
	return &GameSession{
		sessionId:     uuid.NewString(),
		status:        status.SESSION_WAITING,
		host:          host,
		name:          name,
		minPlayer:     minPlayer,
		currentPlayer: uint(0),
		maxPlayer:     maxPlayer,
		wordList:      wordList,
		players:       make(map[string]gameClient.IGameClient),
	}
}

func (gs *GameSession) SetJoinedPlayer(uid string, player gameClient.IGameClient) {
	gs.Lock()
	defer gs.Unlock()
	gs.currentPlayer += 1
	gs.players[uid] = player
}

func (gs *GameSession) SetExitedPlayer(uid string) {
	gs.Lock()
	defer gs.Unlock()

	player, ok := gs.players[uid]
	if ok {
		delete(gs.players, player.GetClientId())
		gs.currentPlayer -= 1
	}

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

func (gs *GameSession) GetSessionHost() gameClient.IGameClient {
	return gs.host
}

func (gs *GameSession) GetSessionName() string {
	return gs.name
}

func (gs *GameSession) GetSessionPlayerInfo() (min uint, max uint, current uint) {
	return gs.minPlayer, gs.maxPlayer, gs.currentPlayer
}

func (gs *GameSession) GetSessionWordList() []string {
	return gs.wordList
}
