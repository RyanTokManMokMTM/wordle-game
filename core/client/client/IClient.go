package client

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/client/internal/types"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
)

type IClient interface {
	Run()
	Close()

	SendToServer(pkgType string, data []byte)

	SetUserId(id string)
	SetUserName(name string)
	SetRenderEvent(eventType string, data *packet.BasicResponseType)
	SetIsWritable()
	SetRenderEventName(eventType string)

	GetUserId() string
	GetUserName() string
	GetInput() chan string
	GetRenderEventName() string
	GetRenderEvent() chan types.RenderEvent
}
