package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameClient/internal/config"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameClient/render"
	"io"
	"log"
	"net"
	"sync"
)

type Client struct {
	sync.Once
	userId      string
	userName    string
	host        string
	port        uint
	networkType string

	conn        net.Conn
	close       chan struct{}
	messageChan chan []byte
	sendChan    chan []byte
	modeChan    chan int
}

func NewClient(c config.Config) IClient {
	return &Client{
		host:        c.Host,
		port:        c.Port,
		networkType: c.NetworkType,
		close:       make(chan struct{}),
		messageChan: make(chan []byte),
		sendChan:    make(chan []byte),
		modeChan:    make(chan int),
	}
}

func (c *Client) Run() {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	source := fmt.Sprintf("%s:%d", c.host, c.port)
	var err error
	c.conn, err = net.Dial(c.networkType, source)
	if err != nil {
		log.Fatal(err)
	}
	go c.onListen(waitGroup)
	go c.read()
	log.Println("Client is connected to ", source)

	c.renderInitPage()

	waitGroup.Wait()

}

func (c *Client) renderInitPage() {
	render.SetUpClientPage(func(name string) {
		c.SetUserName(name)

		//Sending request to update username

		req := packet.EstablishReq{
			PlayerName: name,
		}

		dataBytes, err := serializex.Marshal(req)
		if err != nil {
			log.Println(err)
			return
		}
		go c.SendRequest(packetType.ESTABLISH, dataBytes)
	})
}

func (c *Client) SetUserId(id string) {
	c.userId = id
}

func (c *Client) SetUserName(name string) {
	c.userName = name
}
func (c *Client) GetUserId() string {
	return c.userId
}

func (c *Client) GetUserName() string {
	return c.userName
}

func (c *Client) Close() {
	c.Once.Do(func() {
		c.close <- struct{}{}
		if err := c.conn.Close(); err != nil {
			log.Println(err)
		}
		close(c.close) // Close the channel
		close(c.modeChan)
		close(c.sendChan)
		close(c.messageChan)

	})
}

func (c *Client) read() {
	defer func() {
		c.Close()
	}()

	for {
		var msgLen uint32
		if err := binary.Read(c.conn, binary.BigEndian, &msgLen); err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("binary reading error :", err)
			}
			return
		}

		data := make([]byte, msgLen)
		_, err := c.conn.Read(data)
		if err != nil {
			log.Println("reading error :", err)
			return
		}
		c.messageChan <- data
	}

}

func (c *Client) onListen(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("onListen is ended")
		c.Close()
		wg.Add(-1)
	}()
	for {
		select {
		case <-c.close:
			log.Println("connection closing...")
			return
		case data, ok := <-c.messageChan: //Message from server.
			if !ok {
				log.Fatal("Message channel is closed.")
				return
			}

			var msg *packet.BasicPacket
			if err := serializex.Unmarshal(data, &msg); err != nil {
				log.Print("json err :", err)
				return
			}

			c.handleServerEvent(*msg)
			break

		case dataBytes, ok := <-c.sendChan:
			if !ok {
				log.Print("sending channel is closed")
				return
			}

			log.Println("sending data to server")
			if err := utils.SendMessage(c.conn, dataBytes); err != nil {
				log.Println(err)
			}
			break
		case mode, ok := <-c.modeChan:
			if !ok {
				log.Print("mode channel is closed")
				return
			}

			switch mode {
			case 1:
				fmt.Println("mode is ", mode)
				render.CreateRoomPage(func(roomName string, minPlayer, maxPlayer uint, wordList []string) {
					req := packet.CreateRoomReq{
						UserId:    c.GetUserId(),
						RoomName:  roomName,
						MinPlayer: minPlayer,
						MaxPlayer: maxPlayer,
						WordList:  wordList,
					}

					dataBytes, err := serializex.Marshal(&req)
					if err != nil {
						log.Println(err)
						return
					}

					go c.SendRequest(packetType.CREATE_ROOM, dataBytes)

				})
			case 2:
				fmt.Println("mode is ", mode)
			default:
				fmt.Println("mode is not supported")
			}

			break
		}

	}

}

// Define some function for different event
func (c *Client) handleServerEvent(pk packet.BasicPacket) {
	switch pk.PacketType {
	case packetType.ESTABLISH:
		var resp packet.EstablishResp
		if err := serializex.Unmarshal(pk.Data, &resp); err != nil {
			log.Println(err)
			c.Close()
			return
		}

		c.SetUserId(resp.UserId)
		c.SetUserName(resp.Name)

		log.Println(c.GetUserId())
		log.Println(c.GetUserName())
		go render.MainPage(func(mode int) {
			c.modeChan <- mode
		})
		break
	case packetType.CREATE_ROOM:
		//render create room result page.
		go render.CreateRoomResultPage(pk.Data, c.GetUserId())
		break
	case packetType.JOIN_ROOM:
		//render join room result page.
		go render.JoinRoomResultPage(pk.Data)
		break
	case packetType.EXIT_ROOM:
		//render exit room result page.
		go render.ExistRoomResultPage(pk.Data)
		break
	case packetType.GET_SESSION_INFO:
		//render get session result page.
		break
	default:
		log.Println("Packet not supported")
		return
	}
}

func (c *Client) SendRequest(pkgType string, data []byte) {
	log.Printf("Sending package to server with %s.\n", pkgType)
	switch pkgType {
	case packetType.ESTABLISH:
		fallthrough
	case packetType.CREATE_ROOM:
		fallthrough
	case packetType.JOIN_ROOM:
		fallthrough
	case packetType.EXIT_ROOM:
		fallthrough
	case packetType.GET_SESSION_INFO:
		newPk := packet.NewPacket(pkgType, data)

		dataBytes, err := serializex.Marshal(&newPk)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Sent to send channel")
		c.sendChan <- dataBytes //send to server

	default:
		fmt.Println("Packet type not supported")
	}
}
