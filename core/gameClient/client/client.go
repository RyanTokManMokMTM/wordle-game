package client

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameClient/internal/config"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	host        string
	port        uint
	networkType string

	conn  net.Conn
	close chan struct{}
	input chan struct{}
}

func NewClient(c config.Config) IClient {
	return &Client{
		host:        c.Host,
		port:        c.Port,
		networkType: c.NetworkType,
		close:       make(chan struct{}),
		input:       make(chan struct{}),
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
		if err := c.conn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Start a Wordle Game.")
	fmt.Println("---------------------")

	go c.read()
	go c.write()

	select {
	case <-c.close:
		return
	}

}

func (c *Client) read() {
	for {
		data := make([]byte, 256)
		_, err := c.conn.Read(data)
		if err != nil {
			c.close <- struct{}{}
			return
		}

		fmt.Print(string(data))
	}

}

func (c *Client) write() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			c.close <- struct{}{}
			return
		}
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" {
			c.close <- struct{}{}
			return
		}

		_, err = c.conn.Write([]byte(input))
		if err != nil {
			c.close <- struct{}{}
			return
		}
	}
}
