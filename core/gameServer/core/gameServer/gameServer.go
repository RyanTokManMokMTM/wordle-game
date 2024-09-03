package gameServer

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/manager"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/internal/config"
	"sync"

	"log"
	"net"
)

type GameServer struct {
	sync.Mutex
	host           string
	port           uint
	networkType    string
	listener       net.Listener
	round          uint
	wordList       []string
	sessionManager manager.IGameSessionManager
	clientManager  manager.IGameClientManager

	closedClientChan  chan string //which client is disconnected.
	closedSessionChan chan string //which session is closed

	registerClientChan chan gameClient.IGameClient
}

func NewGameServer(c *config.Config) IGameServer {
	return &GameServer{
		host:               c.Host,
		port:               c.Port,
		networkType:        c.NetworkType,
		round:              c.Round,
		wordList:           c.WordList,
		sessionManager:     manager.NewGameSessionManager(),
		clientManager:      manager.NewGameClientManager(),
		registerClientChan: make(chan gameClient.IGameClient),
		closedClientChan:   make(chan string),
		closedSessionChan:  make(chan string),
	}
}

func (gs *GameServer) Listen() error {
	source := fmt.Sprintf("%s:%d", gs.host, gs.port)
	var err error
	gs.listener, err = net.Listen(gs.networkType, source)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = gs.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go gs.EventListener() //Listening to different channel event

	fmt.Printf("Server listen on %s\n", source)

	for {
		conn, err := gs.listener.Accept()
		fmt.Println("A client has been accepted.")
		if err != nil {
			fmt.Printf("connection accept error %s\n", err.Error())
			continue
		}

		newClient := gameClient.NewGameClient(conn)
		gs.registerClientChan <- newClient

		go newClient.Read()
		go newClient.Write()
		go func() {
			//To get closed signal
			_ = <-newClient.GetClosedEvent()
			userId := newClient.GetClientId()

			fmt.Println("sending closed userId to closedChannel : ", userId)
			gs.closedClientChan <- userId
		}()

	}

}

func (gs *GameServer) EventListener() {
	for {
		select {
		case client, ok := <-gs.registerClientChan:
			if !ok {
				fmt.Println("register-client channel is closed")
				return
			}

			fmt.Println("New client registered")
			gs.clientManager.SetGameClient(client.GetClientId(), client)
			break

		case clientId, ok := <-gs.closedClientChan:
			if !ok {
				fmt.Println("closed-client channel is closed")
				return
			}

			gs.clientManager.RemoveGameClient(clientId)
			break

		case sessionId, ok := <-gs.closedSessionChan:
			if !ok {
				fmt.Println("closed-session channel is closed")
				return
			}
			gs.sessionManager.RemoveGameSession(sessionId)
			break

		}
	}
}

func (gs *GameServer) Close() error {
	fmt.Println("Server is closing.")
	close(gs.closedSessionChan)
	close(gs.closedClientChan)
	close(gs.registerClientChan)
	return gs.listener.Close()
}
