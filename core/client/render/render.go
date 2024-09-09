package render

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/color"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/renderEvent"
	"log"
	"os"
)

type Render struct {
	c client.IClient
}

func NewRender(c client.IClient) *Render {
	return &Render{
		c: c,
	}
}

func (rd *Render) Run() {
	go rd.onListen()
}

func (rd *Render) onListen() {
	for {
		select {
		case event, ok := <-rd.c.GetRenderEvent():
			if !ok {
				return
			}
			switch event.EventType {
			case renderEvent.INIT_PAGE:
				setUpClientPage(rd.c, func(name string) {
					rd.c.SetUserName(name)
					req := packet.EstablishReq{
						PlayerName: name,
					}

					dataBytes, err := serializex.Marshal(req)
					if err != nil {
						log.Fatal(err)
						return
					}
					rd.c.SendToServer(packetType.ESTABLISH, dataBytes)
				})
				break
			case renderEvent.HOME_PAGE:
				mode := mainPage(rd.c)
				switch mode {
				case 1:
					createRoomPage(rd.c, func(roomName string, wordList []string) {
						createRoomReq := packet.CreateRoomReq{
							UserId:   rd.c.GetUserId(),
							RoomName: roomName,
							WordList: wordList,
						}

						dataBytes, err := serializex.Marshal(&createRoomReq)
						if err != nil {
							log.Fatal(err)
							return
						}
						rd.c.SendToServer(packetType.CREATE_ROOM, dataBytes)
					})
				case 2:
					getRoomListReq := packet.GetRoomListInfoReq{
						UserId: rd.c.GetUserId(),
					}

					dataBytes, err := serializex.Marshal(&getRoomListReq)
					if err != nil {
						log.Fatal(err)
						return
					}

					rd.c.SendToServer(packetType.ROOM_LIST_INFO, dataBytes)
				default:
					fmt.Println(color.Red + "Not supported" + color.Reset)
					os.Exit(0)
				}
				break
			case renderEvent.CREATE_ROOM:
				go createRoomResultPage(rd.c, event.Data, func(mode uint, roomId string) {
					switch mode {
					case 0:
						exitRoomReq := packet.ExitRoomReq{
							UserId: rd.c.GetUserId(),
							RoomId: roomId,
						}

						dataBytes, err := serializex.Marshal(&exitRoomReq)
						if err != nil {
							log.Fatal(err)
							return
						}

						rd.c.SendToServer(packetType.EXIT_ROOM, dataBytes)
					case 1:
						gameStartReq := packet.GameStartReq{
							RoomId: roomId,
						}

						dataBytes, err := serializex.Marshal(&gameStartReq)
						if err != nil {
							log.Fatal(err)
							return
						}

						rd.c.SendToServer(packetType.START_GAME, dataBytes)
						break
					default:
						fmt.Println(color.Red + "Not supported" + color.Reset)
						os.Exit(0)
						return
					}
				}, func(roomId, message string) {
					messageReq := packet.GameRoomChatMessage{
						UserId:  rd.c.GetUserId(),
						RoomId:  roomId,
						Message: message,
					}

					dataBytes, err := serializex.Marshal(&messageReq)
					if err != nil {
						log.Fatal(err)
						return
					}

					rd.c.SendToServer(packetType.ROOM_CHAT_MESSAGE, dataBytes)
				})
				break
			case renderEvent.ROOM_LIST_PAGE:
				joinRoomPage(rd.c, event.Data, func(roomId string, isLeave bool) {
					if isLeave {
						go rd.c.SetRenderEvent(renderEvent.HOME_PAGE, nil)
						return
					}

					joinRoomReq := packet.JoinRoomReq{
						UserId: rd.c.GetUserId(),
						RoomId: roomId,
					}

					dataBytes, err := serializex.Marshal(&joinRoomReq)
					if err != nil {
						log.Fatal(err)
						return
					}

					rd.c.SendToServer(packetType.JOIN_ROOM, dataBytes)
				})
				break
			case renderEvent.JOIN_ROOM:
				go joinRoomResultPage(rd.c, event.Data, func(mode uint, roomId string) {
					switch mode {
					case 0:
						exitRoomReq := packet.ExitRoomReq{
							UserId: rd.c.GetUserId(),
							RoomId: roomId,
						}

						dataBytes, err := serializex.Marshal(&exitRoomReq)
						if err != nil {
							log.Fatal(err)
							return
						}

						rd.c.SendToServer(packetType.EXIT_ROOM, dataBytes)
					default:
						fmt.Println("Create room host mode not support")
						os.Exit(0)
						return
					}
				}, func(roomId, message string) {
					messageReq := packet.GameRoomChatMessage{
						UserId:  rd.c.GetUserId(),
						RoomId:  roomId,
						Message: message,
					}

					dataBytes, err := serializex.Marshal(&messageReq)
					if err != nil {
						log.Fatal(err)
						return
					}

					rd.c.SendToServer(packetType.ROOM_CHAT_MESSAGE, dataBytes)
				})
				break
			case renderEvent.START_GAME:
				gameStartingPage()
				break
			case renderEvent.GAME_PAGE:
				gamingOutPut(rd.c, event.Data)
				break
			case renderEvent.GAME_NOTIFICATION:
				notificationOutput(rd.c, event.Data)
				break
			case renderEvent.ENED_GAME:
				go endingGamePage(rd.c, event.Data, func(roomId string) {
					exitRoom := packet.ExitRoomReq{
						UserId: rd.c.GetUserId(),
						RoomId: roomId,
					}

					dataBytes, err := serializex.Marshal(&exitRoom)
					if err != nil {
						log.Fatal(err)
						return
					}

					rd.c.SendToServer(packetType.EXIT_ROOM, dataBytes)
				})
				break
			}
		}
	}
}
