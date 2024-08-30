package client

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/internal/config"
	"math/rand"
)

type Client struct {
	totalRound     uint
	wordList       []string
	guessingWord   string
	remindingRound uint
	wordHistory    []string
}

func NewClient(c config.Config) *Client {
	return &Client{
		totalRound:     c.Round,
		wordList:       c.WordList,
		remindingRound: c.Round,
		wordHistory:    make([]string, 0),
	}
}

func (c *Client) SetWordHistory(w string) {
	c.wordHistory = append(c.wordHistory, w)
}

func (c *Client) SetGuessingWord() {
	maxLen := len(c.GetWordList()) - 1
	index := rand.Intn(maxLen)
	c.guessingWord = c.wordList[index]
}

func (c *Client) IsGameOver() bool {
	return false
}

func (c *Client) Reset() {
	c.guessingWord = ""             //Remove guessing word
	c.remindingRound = c.totalRound //Reset total round
	clear(c.wordHistory)            //Reset Word History
}

func (c *Client) GetTotalRound() uint {
	return c.totalRound
}
func (c *Client) GetWordList() []string {
	return c.wordList
}
func (c *Client) GetGuessingWord() string {
	return c.guessingWord
}
func (c *Client) GetWordHistory() []string {
	return c.wordHistory
}

func (c *Client) GetRemindingRound() uint {
	return c.remindingRound
}
