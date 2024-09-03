package gameClient

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
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

	readChannel  chan []byte
	writeChannel chan []byte

	isClosed chan struct{}
}

func NewGameClient(conn net.Conn) IGameClient {
	return &GameClient{
		clientId:     uuid.NewString(),
		conn:         conn,
		readChannel:  make(chan []byte),
		writeChannel: make(chan []byte),
		isClosed:     make(chan struct{}),
	}
}

func (gc *GameClient) Read() {
	defer func() {
		gc.Closed()
	}()

	for {
		reader := bufio.NewReader(gc.conn)
		data := make([]byte, 256)
		n, err := reader.Read(data[:])
		if err != nil {
			log.Println(err)
			return
		}

		gc.readChannel <- data[:n]
	}
}

func (gc *GameClient) Closed() {
	gc.Do(func() {
		fmt.Println("Client closed")
		gc.isClosed <- struct{}{}
		gc.conn.Close()

		close(gc.readChannel)
		close(gc.writeChannel)
		close(gc.isClosed)
	})
}

func (gc *GameClient) Write() {
	defer func() {
		gc.Closed()
	}()

	//Write any message to this client
	for {
		msg, ok := <-gc.writeChannel
		if !ok {
			log.Println("Write channel is closed.")
			return
		}

		//write a message to this client
		if err := writeMessage(gc.conn, msg); err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println(err)
				return
			}
		}
	}
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

func writeMessage(conn net.Conn, dataBytes []byte) error {
	msgLen := uint32(len(dataBytes))
	err := binary.Write(conn, binary.BigEndian, msgLen)
	if err != nil {
		return err
	}

	_, err = conn.Write(dataBytes)
	if err != nil {
		return err
	}

	return nil
}
