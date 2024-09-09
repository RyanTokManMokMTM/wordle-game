package types

import "github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"

type RenderEvent struct {
	Data      *packet.BasicResponseType
	EventType string
}

func NewRenderEvent(eventType string, resp *packet.BasicResponseType) RenderEvent {
	return RenderEvent{
		Data:      resp,
		EventType: eventType,
	}
}
