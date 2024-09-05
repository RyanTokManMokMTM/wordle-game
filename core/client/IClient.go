package client

type IClient interface {
	// GetTotalRound get total round of the game
	GetTotalRound() uint

	// GetWordList get the given word list
	GetWordList() []string
}
