package client

type IClient interface {
	SetWordHistory(string)
	SetGuessingWord()

	GetTotalRound() uint
	GetWordList() []string
	GetGuessingWord() string
	GetWordHistory() []string

	Reset()
}
