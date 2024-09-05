package manager

import "github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameRoom"

type IGameRoomManager interface {
	SetGameRoom(roomId string, room gameRoom.IGameRoom)
	RemoveGameRoom(roomId string)

	GetGameRoom(roomId string) (gameRoom.IGameRoom, bool)
	GetGameRoomList() []gameRoom.IGameRoom
}
