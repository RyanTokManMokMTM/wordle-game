package gameClient

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/internal/logic"
	"math/rand"
	"net"
)

type GameClient struct {
	conn net.Conn

	totalRound   uint
	wordList     []string
	guessingWord string
	wordHistory  []string
}

func NewGameClient(conn net.Conn, round uint, wordList []string) IGameClient {
	return &GameClient{
		conn:        conn,
		totalRound:  round,
		wordList:    wordList,
		wordHistory: make([]string, 0),
	}
}

func (gc *GameClient) HandleRequest() {
	gc.SetGuessingWord()
	logic.GameLogic(gc.guessingWord, gc.totalRound, gc.conn)
	fmt.Println("Bye")
	gc.conn.Close()
}

func (gc *GameClient) SetWordHistory(w string) {
	gc.wordHistory = append(gc.wordHistory, w)
}

func (gc *GameClient) SetGuessingWord() {
	rand.NewSource(0)
	maxLen := len(gc.GetWordList()) - 1
	index := rand.Intn(maxLen)
	gc.guessingWord = gc.wordList[index]
}

func (gc *GameClient) Reset() {
	gc.guessingWord = ""  //Remove guessing word
	clear(gc.wordHistory) //Reset Word History
}

func (gc *GameClient) GetTotalRound() uint {
	return gc.totalRound
}

func (gc *GameClient) GetWordList() []string {
	return gc.wordList
}

func (gc *GameClient) GetGuessingWord() string {
	return gc.guessingWord
}

func (gc *GameClient) GetWordHistory() []string {
	return gc.wordHistory
}

func (gc *GameClient) GetConn() net.Conn {
	return gc.conn
}
