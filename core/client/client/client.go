package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/internal/config"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/internal/types"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/renderEvent"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	sync.Once
	userId      string
	userName    string
	host        string
	port        uint
	networkType string

	conn            net.Conn
	close           chan struct{}
	messageChan     chan []byte
	sendChan        chan []byte
	inputChan       chan string
	isWritable      chan struct{}
	renderEventChan chan types.RenderEvent

	previousRenderEvent string
	currentRenderEvent  string
}

func NewClient(c config.Config) IClient {
	return &Client{
		host:            c.Host,
		port:            c.Port,
		networkType:     c.NetworkType,
		close:           make(chan struct{}),
		messageChan:     make(chan []byte),
		sendChan:        make(chan []byte),
		isWritable:      make(chan struct{}, 2), //at most 10 uses?
		renderEventChan: make(chan types.RenderEvent),
		inputChan:       make(chan string),
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
	go c.write()
	log.Println("Client is connected to ", source)

	c.SetRenderEvent(renderEvent.INIT_PAGE, nil)
	waitGroup.Wait()

}

func (c *Client) SetRenderEventName(eventType string) {
	c.previousRenderEvent = c.currentRenderEvent
	c.currentRenderEvent = eventType
}

func (c *Client) SetRenderEvent(eventType string, data []byte) {
	c.renderEventChan <- types.NewRenderEvent(eventType, data)
}

func (c *Client) GetRenderEventName() string {
	return c.currentRenderEvent
}

func (c *Client) GetInput() chan string {
	c.SetIsWritable()
	return c.inputChan
}

func (c *Client) GetRenderEvent() chan types.RenderEvent {
	return c.renderEventChan
}

func (c *Client) SetUserId(id string) {
	c.userId = id
}

func (c *Client) SetUserName(name string) {
	c.userName = name
}

func (c *Client) SetIsWritable() {
	c.isWritable <- struct{}{}
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
		//close(c.modeChan)
		close(c.sendChan)
		close(c.messageChan)
		close(c.renderEventChan)
		close(c.inputChan)

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

func (c *Client) write() {
	defer func() {
		c.Close()
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case _, ok := <-c.isWritable:
			if !ok {
				return
			}
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Println("reading input error :", err)
				return
			}
			inputInfo := strings.Trim(input, "\r\n")

			eventName := c.GetRenderEventName()

			switch eventName {
			case renderEvent.GAME_PAGE:

				req := packet.PlayingGameReq{
					Input: []byte(inputInfo),
				}

				dataBytes, err := serializex.Marshal(&req)
				if err != nil {
					log.Fatal(err)
					return
				}

				pk := packet.NewPacket(packetType.PLAYING_GAME, dataBytes)
				pkData, err := serializex.Marshal(&pk)
				if err != nil {
					log.Fatal(err)
					return
				}
				c.sendChan <- pkData
				break
			default:
				c.inputChan <- inputInfo
				break
			}
		}
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
			if err := utils.SendMessage(c.conn, dataBytes); err != nil {
				log.Println(err)
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
		c.SetRenderEventName(renderEvent.HOME_PAGE)
		c.SetRenderEvent(renderEvent.HOME_PAGE, nil)
		break
	case packetType.ROOM_LIST_INFO:
		c.SetRenderEventName(renderEvent.ROOM_LIST_PAGE)
		c.SetRenderEvent(renderEvent.ROOM_LIST_PAGE, pk.Data)
		break
	case packetType.CREATE_ROOM:
		//render create room result page.
		c.SetRenderEventName(renderEvent.CREATE_ROOM)
		c.SetRenderEvent(renderEvent.CREATE_ROOM, pk.Data)
		break
	case packetType.JOIN_ROOM:
		//render join room result page.
		c.SetRenderEventName(renderEvent.JOIN_ROOM)
		c.SetRenderEvent(renderEvent.JOIN_ROOM, pk.Data)
		break
	case packetType.EXIT_ROOM:
		c.SetRenderEventName(renderEvent.HOME_PAGE)
		c.SetRenderEvent(renderEvent.HOME_PAGE, nil)
		break
	case packetType.START_GAME:
		c.SetRenderEvent(renderEvent.START_GAME, pk.Data)
		break
	case packetType.PLAYING_GAME:
		if strings.Compare(c.previousRenderEvent, renderEvent.JOIN_ROOM) == 0 {
			c.inputChan <- "-1" //To disable join room event stdin
		}
		c.SetRenderEventName(renderEvent.GAME_PAGE)
		c.SetRenderEvent(renderEvent.GAME_PAGE, pk.Data)
		break
	case packetType.GAME_NOTIFICATION:
		c.SetRenderEvent(renderEvent.GAME_NOTIFICATION, pk.Data)
		break
	case packetType.END_GAME:
		c.SetRenderEventName(renderEvent.ENED_GAME)
		c.SetRenderEvent(renderEvent.ENED_GAME, pk.Data)
		break
	default:
		log.Println("Packet not supported")
		return
	}
}

func (c *Client) SendToServer(pkgType string, data []byte) {
	switch pkgType {
	case packetType.ESTABLISH:
		fallthrough
	case packetType.CREATE_ROOM:
		fallthrough
	case packetType.ROOM_LIST_INFO:
		fallthrough
	case packetType.JOIN_ROOM:
		fallthrough
	case packetType.START_GAME:
		fallthrough
	case packetType.EXIT_ROOM:
		pk := packet.NewPacket(pkgType, data)

		dataBytes, err := serializex.Marshal(&pk)
		if err != nil {
			log.Println(err)
			return
		}
		c.sendChan <- dataBytes //send to server
	default:
		fmt.Println("Packet type not supported")
	}
}
