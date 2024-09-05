package types

type RenderEvent struct {
	Data      []byte
	EventType string
}

func NewRenderEvent(eventType string, data []byte) RenderEvent {
	return RenderEvent{
		Data:      data,
		EventType: eventType,
	}
}
