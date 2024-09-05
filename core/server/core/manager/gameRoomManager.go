package manager

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameRoom"
	"sync"
)

type GameRoomManager struct {
	sync.Mutex
	rooms map[string]gameRoom.IGameRoom
}

func NewGameRoomManager() IGameRoomManager {
	return &GameRoomManager{
		rooms: make(map[string]gameRoom.IGameRoom),
	}
}

func (gsm *GameRoomManager) SetGameRoom(roomId string, room gameRoom.IGameRoom) {
	gsm.Lock()
	defer gsm.Unlock()
	_, ok := gsm.rooms[roomId]
	if !ok {
		gsm.rooms[roomId] = room
	}
}
func (gsm *GameRoomManager) RemoveGameRoom(roomId string) {
	s, ok := gsm.rooms[roomId]
	if ok {
		delete(gsm.rooms, roomId)
		s.Close()
	}
}
func (gsm *GameRoomManager) GetGameRoom(roomId string) (gameRoom.IGameRoom, bool) {
	s, ok := gsm.rooms[roomId]
	if !ok {
		return nil, false
	}
	return s, true
}

func (gsm *GameRoomManager) GetGameRoomList() []gameRoom.IGameRoom {
	roomList := make([]gameRoom.IGameRoom, 0)
	for _, s := range gsm.rooms {
		roomList = append(roomList, s)
	}

	return roomList
}
