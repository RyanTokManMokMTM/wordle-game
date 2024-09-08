package packet

import "github.com/google/uuid"

type BasicPacket struct {
	PacketId   string `json:"packet_id"`
	PacketType string `json:"packet_type"`
	Data       []byte `json:"data"` //Decode by to the packet type
}

func NewPacket(pkType string, data []byte) *BasicPacket {
	return &BasicPacket{
		PacketId:   uuid.NewString(),
		PacketType: pkType,
		Data:       data,
	}
}

type EstablishReq struct {
	PlayerName string `json:"user_name"`
}

type EstablishResp struct {
	UserId string `json:"user_id"`
	Name   string `json:"user_name"`
}

type CreateRoomReq struct {
	UserId   string   `json:"user_id"`
	RoomName string   `json:"room_name"`
	WordList []string `json:"word_list"`
}

type CreateRoomResp struct {
	Code uint `json:"code"`
	GameRoomInfoPacket
}

type JoinRoomReq struct {
	UserId string `json:"user_id"`
	RoomId string `json:"room_id"`
}

type JoinRoomResp struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	GameRoomInfoPacket
}

type ExitRoomReq struct {
	UserId string `json:"user_id"`
	RoomId string `json:"room_id"`
}

type ExitRoomResp struct {
	Code uint `json:"code"`
}

type GetRoomListInfoReq struct {
	UserId string `json:"user_id"`
}

type GetRoomListInfoResp struct {
	Code  uint                 `json:"code"`
	Rooms []GameRoomInfoPacket `json:"rooms"`
}

type GetRoomInfoReq struct {
	UserId string `json:"user_id"`
	RoomId string `json:"room_id"`
}

type GetRoomInfoResp struct {
	Code uint `json:"code"`
	GameRoomInfoPacket
}

type GameRoomInfoPacket struct {
	RoomId       string `json:"room_id"`
	RoomName     string `json:"room_name"`
	RoomHostName string `json:"room_host_name"`
	RoomHostId   string `json:"room_host_id"`

	RoomStatus string `json:"room_status"`
}

type GameStartReq struct {
	RoomId string `json:"room_id"`
}

type PlayingGameReq struct {
	Input []byte `json:"input"`
}

type PlayingGameResp struct {
	OutputColorASNI string `json:"output_color_asni"`
	IsWritable      bool   `json:"is_writable"` //Indicate user need to performance an input action
	GameMessage     []byte `json:"game_message"`
}

type NotifyPlayer struct {
	Type    string `json:"type"`
	Message []byte `json:"message"`
}

type EndingGameResp struct {
	OutputColorASNI string `json:"output_color_asni"`
	RoomId          string `json:"room_id"`
	Message         []byte `json:"message"`
}

type GameRoomChatMessage struct {
	UserId  string `json:"user_id"`
	RoomId  string `json:"roomId"`
	Message string `json:"message"`
}
