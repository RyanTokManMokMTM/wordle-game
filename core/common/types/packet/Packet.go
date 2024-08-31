package packet

import "github.com/google/uuid"

type Packet struct {
	PacketId   string `json:"packet_id"`
	PacketType string `json:"packet_type"`
	GamePacket
}

type GamePacket struct {
	IsWritable  bool   `json:"is_writable"` //Indicate user need to performance an input action
	GameMessage []byte `json:"game_message"`
}

func NewPacket(msgType string, IsWritable bool, gameMessage string) Packet {
	return Packet{
		PacketId:   uuid.NewString(),
		PacketType: msgType,
		GamePacket: GamePacket{
			IsWritable:  IsWritable,
			GameMessage: []byte(gameMessage),
		},
	}
}
