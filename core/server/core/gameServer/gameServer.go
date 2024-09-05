package gameServer

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameClient"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/internal/config"

	"log"
	"net"
)

type GameServer struct {
	host        string
	port        uint
	networkType string
	listener    net.Listener
	round       uint
	wordList    []string
}

func NewGameServer(c config.Config) IGameServer {
	return &GameServer{
		host:        c.Host,
		port:        c.Port,
		networkType: c.NetworkType,
		round:       c.Round,
		wordList:    c.WordList,
	}
}

// Listen starting a server
func (gs *GameServer) Listen() error {
	source := fmt.Sprintf("%s:%d", gs.host, gs.port)
	var err error
	gs.listener, err = net.Listen(gs.networkType, source)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = gs.listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Printf("Server listen on %s\n", source)
	for {
		conn, err := gs.listener.Accept()
		fmt.Println("A client has been accepted.")
		if err != nil {
			fmt.Printf("connection accept error %s\n", err.Error())
			continue
		}

		newClient := gameClient.NewGameClient(conn, gs.round, gs.wordList)
		go newClient.HandleRequest()
	}

}

func (gs *GameServer) Close() error {
	fmt.Println("Server is closing.")
	return gs.listener.Close()
}
