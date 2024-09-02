package client

type IClient interface {
	GetTotalRound() uint
	GetWordList() []string
}
