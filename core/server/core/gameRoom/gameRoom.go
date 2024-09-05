package gameRoom

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/status"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gamePlayer"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/internal/logic"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"strings"
	"sync"
)

type GameRoom struct {
	sync.Mutex
	roomId   string
	host     gamePlayer.IGamePlayer
	name     string
	wordList []string
	status   string
	players  map[string]gamePlayer.IGamePlayer

	totalRound          uint
	currentGuessingWord string

	isFinishedGame chan struct{}               //received a signal that the game is finished(all player)
	gameOverPlayer chan gamePlayer.IGamePlayer //received player who is finished
	endedPlayer    []string
}

func NewGameRoom(host gamePlayer.IGamePlayer, name string, wordList []string, totalRound uint) IGameRoom {
	return &GameRoom{
		roomId: uuid.NewString(),
		status: status.ROOM_STAUS_WAITING,
		host:   host,
		name:   name,

		wordList:   wordList,
		players:    make(map[string]gamePlayer.IGamePlayer),
		totalRound: totalRound,

		isFinishedGame: make(chan struct{}),
		gameOverPlayer: make(chan gamePlayer.IGamePlayer),
		endedPlayer:    make([]string, 0),
	}
}

func (gr *GameRoom) AddPlayer(uid string, p gamePlayer.IGamePlayer) {
	gr.Lock()
	defer gr.Unlock()
	gr.players[uid] = p
}

func (gr *GameRoom) RemovePlayer(uid string) bool {
	gr.Lock()
	defer gr.Unlock()

	p, ok := gr.players[uid]
	if ok {
		delete(gr.players, p.GetClient().GetClientId())
		return true
	}
	return false

}

func (gr *GameRoom) SetRoomStatus(status string) {
	gr.status = status
}

func (gr *GameRoom) SetGuessingWord() {
	maxLen := len(gr.GetRoomWordList()) - 1
	if maxLen < 0 {
		log.Fatal("word list is empty")
		return
	}

	if maxLen == 0 {
		gr.currentGuessingWord = gr.GetRoomWordList()[0]
		return
	}

	rand.NewSource(0)
	index := rand.Intn(maxLen)
	gr.currentGuessingWord = gr.GetRoomWordList()[index]

}

func (gr *GameRoom) GetAllPlayer() []gamePlayer.IGamePlayer {
	allPlayer := make([]gamePlayer.IGamePlayer, 0)
	for _, p := range gr.players {
		allPlayer = append(allPlayer, p)
	}
	return allPlayer
}

func (gr *GameRoom) GetRoomStatus() string {
	return gr.status
}

func (gr *GameRoom) GetRoomId() string {
	return gr.roomId
}

func (gr *GameRoom) GetRoomHost() gamePlayer.IGamePlayer {
	return gr.host
}

func (gr *GameRoom) GetRoomName() string {
	return gr.name
}

func (gr *GameRoom) GetRoomWordList() []string {
	return gr.wordList
}

func (gr *GameRoom) RemoveAllPlayer() {
	gr.Lock()
	defer gr.Unlock()

	for uid, _ := range gr.players {
		delete(gr.players, uid)
	}
}

// StartGame  staring the game process for this player
func (gr *GameRoom) StartGame(player gamePlayer.IGamePlayer) {
	log.Println("Game stared with player : ", player.GetClient().GetName())
	logic.GameLogic(gr.currentGuessingWord, gr.totalRound, player)

	gr.updateEndedPlayer(player.GetClient().GetClientId())
	if !gr.isOver() {
		gr.notifyPlayer(player)
	} else {
		gr.isFinishedGame <- struct{}{}
	}
}

func (gr *GameRoom) updateEndedPlayer(id string) {
	gr.Lock()
	defer gr.Unlock()
	gr.endedPlayer = append(gr.endedPlayer, id)
}

// isOver is all player done the game
func (gr *GameRoom) isOver() bool {
	size := len(gr.players)
	doneSize := len(gr.endedPlayer)

	return size == doneSize
}

// GetTheGameIsOver getting isFinishedGame signal from channel
func (gr *GameRoom) GetTheGameIsOver() chan struct{} {
	return gr.isFinishedGame
}

func (gr *GameRoom) Close() {
	close(gr.isFinishedGame)
	close(gr.gameOverPlayer)
}

// notifyPlayer with a defined message
func (gr *GameRoom) notifyPlayer(withoutClient gamePlayer.IGamePlayer) {
	players := gr.players
	for _, p := range players {
		var message string
		if strings.Compare(p.GetClient().GetClientId(), withoutClient.GetClient().GetClientId()) == 0 {
			message = fmt.Sprintf("[SYSTEM] You finished the game, please waiting for other player to finish.\n")
		} else {
			message = fmt.Sprintf("[SYSTEM] Player %s is finished the game! His score is %d.\n", withoutClient.GetClient().GetName(), withoutClient.GetScore())
		}
		gr.NotifyPlayerWithMessage(p, message)
	}
}

func (gr *GameRoom) NotifyPlayerWithMessage(p gamePlayer.IGamePlayer, message string) {
	req := packet.NotifyPlayer{
		Message: []byte(message),
	}

	dataBytes, err := serializex.Marshal(&req)
	if err != nil {
		log.Println(err)
		return
	}

	pk := packet.NewPacket(packetType.GAME_NOTIFICATION, dataBytes)
	dataBytes, err = serializex.Marshal(&pk)
	if err != nil {
		log.Println(err)
		return
	}
	if err := utils.SendMessage(p.GetClient().GetConn(), dataBytes); err != nil {
		log.Println(err)
	}

}
