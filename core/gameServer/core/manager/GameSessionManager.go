package manager

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameSession"
	"sync"
)

type GameSessionManager struct {
	sync.Mutex
	sessions map[string]gameSession.IGameSession
}

func NewGameSessionManager() IGameSessionManager {
	return &GameSessionManager{
		sessions: make(map[string]gameSession.IGameSession),
	}
}

func (gsm *GameSessionManager) SetGameSession(sessionId string, session gameSession.IGameSession) {
	gsm.Lock()
	defer gsm.Unlock()

	//TODO: delete the existing session
	_, ok := gsm.sessions[sessionId]
	if !ok {
		gsm.sessions[sessionId] = session
	}
}
func (gsm *GameSessionManager) RemoveGameSession(sessionId string) {
	_, ok := gsm.sessions[sessionId]
	if ok {
		//TODO: What to do to the session?
		delete(gsm.sessions, sessionId)
	}
}
func (gsm *GameSessionManager) GetGameSession(sessionId string) (*gameSession.IGameSession, bool) {
	s, ok := gsm.sessions[sessionId]
	if !ok {
		return nil, false
	}
	return &s, true
}
func (gsm *GameSessionManager) GetGameSessionList() []gameSession.IGameSession {
	sessionList := make([]gameSession.IGameSession, 0)
	for _, s := range gsm.sessions {
		sessionList = append(sessionList, s)
	}

	return sessionList
}
