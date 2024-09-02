package client

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/internal/config"
)

type Client struct {
	totalRound uint
	wordList   []string
}

func NewClient(c config.Config) IClient {
	return &Client{
		totalRound: c.Round,
		wordList:   c.WordList,
	}
}

func (c *Client) GetTotalRound() uint {
	return c.totalRound
}
func (c *Client) GetWordList() []string {
	return c.wordList
}
