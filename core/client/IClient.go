package client

type IClient interface {
	// SetGuessingWord get a guessing word randomly from given list
	SetGuessingWord()

	// GetTotalRound How many total round of the game
	GetTotalRound() uint

	// GetWordList Get the predefined word list
	GetWordList() []string

	// GetGuessingWord Get current guessing word
	GetGuessingWord() string
}
