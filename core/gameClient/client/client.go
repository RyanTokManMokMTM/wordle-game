package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameClient/internal/config"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	sync.Once
	host        string
	port        uint
	networkType string

	conn        net.Conn
	close       chan struct{}
	writable    chan struct{}
	messageChan chan []byte
}

func NewClient(c config.Config) IClient {
	return &Client{
		host:        c.Host,
		port:        c.Port,
		networkType: c.NetworkType,
		close:       make(chan struct{}),
		writable:    make(chan struct{}), //Be able to write
		messageChan: make(chan []byte),
	}
}

func (c *Client) Run() {
	source := fmt.Sprintf("%s:%d", c.host, c.port)
	var err error
	c.conn, err = net.Dial(c.networkType, source)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		close(c.close)    // Close the channel
		close(c.writable) //Close the channel
	}()

	fmt.Println("Start a Wordle Game.")
	fmt.Println("---------------------")

	go c.read()
	go c.write()

	c.onListen()

}
func (c *Client) onListen() {
	for {
		select {
		case <-c.close:
			log.Println("connection closing...")
			return
		case data, ok := <-c.messageChan:
			if !ok {
				log.Fatal("Message channel is closed.")
				return
			}
			var msg *packet.Packet
			if err := serializex.Unmarshal(data, &msg); err != nil {
				log.Print("json err :", err)
				return
			}

			c.handleServerEvent(*msg)
			break
		}
	}

}

func (c *Client) Close() {
	c.Once.Do(func() {
		//TODO: to run it once only
		c.close <- struct{}{}
		if err := c.conn.Close(); err != nil {
			log.Println(err)
		}
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
	//TODO: allowing input from user when receiving an write signed?
	defer func() {
		c.Close()
	}()

	for {
		_, ok := <-c.writable
		if !ok {
			return
		}
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("reading input error :", err)
			return
		}
		inputInfo := strings.Trim(input, "\r\n")

		_, err = c.conn.Write([]byte(inputInfo))
		if err != nil {
			log.Println("write input error :", err)
			return
		}
	}
}

// Define some function for different event
func (c *Client) handleServerEvent(serverMsg packet.Packet) {
	switch serverMsg.PacketType {
	case packetType.IN_GAME:
		fmt.Print(string(serverMsg.GameMessage))
		if serverMsg.IsWritable {
			c.writable <- struct{}{} // be able to input...
		}
	default:
		log.Println("Packet not supported")
		return
	}
}
