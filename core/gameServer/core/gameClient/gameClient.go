package gameClient

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"sync"
)

type GameClient struct {
	sync.Once

	clientId string
	conn     net.Conn
	name     string

	readChan    chan []byte
	writeChan   chan []byte
	messageChan chan packet.BasicPacket

	isClosed chan struct{}
}

func NewGameClient(conn net.Conn) IGameClient {
	return &GameClient{
		clientId:    uuid.NewString(),
		conn:        conn,
		readChan:    make(chan []byte),
		writeChan:   make(chan []byte),
		isClosed:    make(chan struct{}),
		messageChan: make(chan packet.BasicPacket),
	}
}

func (gc *GameClient) Run() {
	go gc.read()
	go gc.write()
	go gc.eventListener()
}

func (gc *GameClient) read() {
	defer func() {
		gc.Closed()
	}()

	for {
		reader := bufio.NewReader(gc.conn)
		var msgLen uint32
		if err := binary.Read(gc.conn, binary.BigEndian, &msgLen); err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("binary reading error :", err)
			}
			return
		}

		data := make([]byte, msgLen)
		_, err := reader.Read(data)
		if err != nil {
			log.Println(err)
			return
		}
		gc.readChan <- data
	}
}

func (gc *GameClient) write() {
	defer func() {
		gc.Closed()
	}()

	//Write any message to this client
	for {
		msg, ok := <-gc.writeChan
		if !ok {
			log.Println("Client's write channel is closed.")
			return
		}

		//write a message to this client
		if err := utils.SendMessage(gc.conn, msg); err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println(err)
				return
			}
		}
	}
}

func (gc *GameClient) eventListener() {
	for {
		select {
		case data, ok := <-gc.readChan:
			if !ok {
				fmt.Println("Client's read channel closed")
				return
			}
			var packetData packet.BasicPacket
			if err := serializex.Unmarshal(data, &packetData); err != nil {
				log.Println("packet serialized error : ", err)
				continue
			}
			log.Printf("%+v", packetData)
			pkType := packetData.PacketType
			pkData := packetData.Data

			switch pkType {
			case packetType.ESTABLISH:
				//TODO: update client information.
				var establishReq packet.EstablishReq
				if err := serializex.Unmarshal(pkData, &establishReq); err != nil {
					log.Println("serialized create room req error : ", err)
					return
				}
				gc.SetUserName(establishReq.PlayerName)
				log.Println("Set client name,", gc.GetName())

				establishResp := packet.EstablishResp{
					UserId: gc.GetClientId(),
					Name:   gc.GetName(),
				}
				dataBytes, err := serializex.Marshal(establishResp)
				if err != nil {
					log.Println(err)
					gc.Closed()
					return
				}

				pkResp := packet.NewPacket(packetType.ESTABLISH, dataBytes)
				dataBytes, err = serializex.Marshal(&pkResp)
				if err != nil {
					log.Println(err)
					gc.Closed()
					return
				}

				if err := utils.SendMessage(gc.GetConn(), dataBytes); err != nil {
					log.Println(err)
					gc.Closed()
					return
				}

				log.Println(gc.GetName())
				log.Println("Client set and connected.")
				break
			default:
				//Handle it by server.
				log.Println("other type?")
				gc.messageChan <- packetData
			}

			break
		}
	}
}

func (gc *GameClient) Closed() {
	gc.Do(func() {
		fmt.Println("Client closed")
		gc.isClosed <- struct{}{}
		if err := gc.conn.Close(); err != nil {
			log.Println(err)
		}

		close(gc.readChan)
		close(gc.writeChan)
		close(gc.isClosed)
		close(gc.messageChan)
	})
}

func (gc *GameClient) SetUserName(name string) {
	gc.name = name
}

func (gc *GameClient) GetConn() net.Conn {
	return gc.conn
}

func (gc *GameClient) GetClientId() string {
	return gc.clientId
}

func (gc *GameClient) GetUserName() string {
	return gc.name
}

func (gc *GameClient) GetClosedEvent() chan struct{} {
	return gc.isClosed
}

func (gc *GameClient) GetMessage() chan packet.BasicPacket {
	return gc.messageChan
}

func (gc *GameClient) GetName() string {
	return gc.name
}
