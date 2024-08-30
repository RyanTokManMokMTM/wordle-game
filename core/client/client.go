package client

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/internal/config"
	"math/rand"
)

type Client struct {
	totalRound   uint
	wordList     []string
	guessingWord string
	wordHistory  []string
}

func NewClient(c config.Config) IClient {
	return &Client{
		totalRound:  c.Round,
		wordList:    c.WordList,
		wordHistory: make([]string, 0),
	}
}

func (c *Client) SetWordHistory(w string) {
	c.wordHistory = append(c.wordHistory, w)
}

func (c *Client) SetGuessingWord() {
	rand.NewSource(0)
	maxLen := len(c.GetWordList()) - 1
	index := rand.Intn(maxLen)
	c.guessingWord = c.wordList[index]
}

func (c *Client) Reset() {
	c.guessingWord = ""  //Remove guessing word
	clear(c.wordHistory) //Reset Word History
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
