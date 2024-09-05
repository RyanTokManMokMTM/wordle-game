package manager

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameClient"
	"sync"
)

type GameClientManager struct {
	sync.Mutex
	clients map[string]gameClient.IGameClient
}

func NewGameClientManager() IGameClientManager {
	return &GameClientManager{
		clients: make(map[string]gameClient.IGameClient),
	}
}

func (gcm *GameClientManager) SetGameClient(clientId string, client gameClient.IGameClient) {
	gcm.Lock()
	defer gcm.Unlock()

	_, ok := gcm.clients[clientId]
	if !ok {
		gcm.clients[clientId] = client
	}
}
func (gcm *GameClientManager) RemoveGameClient(clientId string) {
	_, ok := gcm.clients[clientId]
	if ok {
		delete(gcm.clients, clientId)
	}
}

func (gcm *GameClientManager) GetGameClient(clientId string) (gameClient.IGameClient, bool) {
	s, ok := gcm.clients[clientId]
	if !ok {
		return nil, false
	}
	return s, true
}

func (gcm *GameClientManager) GetGameClientList() []gameClient.IGameClient {
	roomList := make([]gameClient.IGameClient, 0)
	for _, s := range gcm.clients {
		roomList = append(roomList, s)
	}

	return roomList
}
