package client

import "github.com/RyanTokManMokMTM/wordle-game/core/client/internal/types"

type IClient interface {
	Run()
	Close()

	SendRequest(pkgType string, data []byte)

	SetUserId(id string)
	SetUserName(name string)
	SetRenderEvent(eventType string, data []byte)
	SetIsWritable()
	SetRenderEventName(eventType string)

	GetUserId() string
	GetUserName() string
	GetInput() chan string
	GetRenderEventName() string
	GetRenderEvent() chan types.RenderEvent
}
