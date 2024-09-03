package gameServer

type IGameServer interface {
	Listen() error
	Close() error

	EventListener()
}
