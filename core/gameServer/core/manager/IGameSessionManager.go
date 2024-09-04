package manager

import "github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameSession"

type IGameSessionManager interface {
	SetGameSession(sessionId string, session gameSession.IGameSession)
	RemoveGameSession(sessionId string)

	GetGameSession(sessionId string) (gameSession.IGameSession, bool)
	GetGameSessionList() []gameSession.IGameSession
}
