package gameServer

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/code"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	utils "github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameClient"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/core/gameSession"
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

	closedClientChan   chan []byte //which client is disconnected.
	closedSessionChan  chan []byte //which session is closed
	messageChan        chan packet.BasicPacket
	registerClientChan chan gameClient.IGameClient
}

func NewGameServer(c *config.Config) IGameServer {
	return &GameServer{
		host:           c.Host,
		port:           c.Port,
		networkType:    c.NetworkType,
		round:          c.Round,
		wordList:       c.WordList,
		sessionManager: manager.NewGameSessionManager(),
		clientManager:  manager.NewGameClientManager(),

		registerClientChan: make(chan gameClient.IGameClient),
		messageChan:        make(chan packet.BasicPacket),
		closedClientChan:   make(chan []byte),
		closedSessionChan:  make(chan []byte),
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

	go gs.eventListener() //Listening to different channel event

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

		go newClient.Run()
		go func() {
			//To get closed signal
			_ = <-newClient.GetClosedEvent()
			userId := newClient.GetClientId()

			gs.closedClientChan <- []byte(userId)
		}()

		go func() {
			//To received message from client
			for {
				msg, ok := <-newClient.GetMessage()
				if !ok {
					log.Printf("client %s message channal closed\n", newClient.GetClientId())
					return
				}

				gs.messageChan <- msg
			}
		}()

	}

}

func (gs *GameServer) eventListener() {
	for {
		select {
		case client, ok := <-gs.registerClientChan:
			if !ok {
				fmt.Println("register-client channel is closed")
				return
			}

			fmt.Println("New client registered")
			gs.clientManager.SetGameClient(client.GetClientId(), client)

			fmt.Println("Current user : ", len(gs.clientManager.GetGameClientList()))
			break

		case clientId, ok := <-gs.closedClientChan:
			if !ok {
				fmt.Println("closed-client channel is closed")
				return
			}
			fmt.Println("client disconnected  : ", string(clientId))
			gs.clientManager.RemoveGameClient(string(clientId))
			break

		case sessionId, ok := <-gs.closedSessionChan:
			if !ok {
				fmt.Println("closed-session channel is closed")
				return
			}
			gs.sessionManager.RemoveGameSession(string(sessionId))
			break

		case message, ok := <-gs.messageChan:
			if !ok {
				fmt.Println("closed-session channel is closed")
				return
			}
			gs.handleMessage(message)
			break
		}
	}
}

func (gs *GameServer) handleMessage(pk packet.BasicPacket) {
	pkType := pk.PacketType
	pkData := pk.Data
	log.Println(pkType)
	switch pkType {
	case packetType.CREATE_ROOM:
		//Create a new room by user
		log.Println("Received an message of create room")
		var createRoomReq packet.CreateRoomReq
		if err := serializex.Unmarshal(pkData, &createRoomReq); err != nil {
			log.Println("serialized create room req error : ", err)
			return
		}

		log.Printf("%+v", createRoomReq)

		//MARK: Get User Client by id
		hostUser, ok := gs.clientManager.GetGameClient(createRoomReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", createRoomReq.UserId)
			return
		}

		//MARK: Creating a session
		newSession := gameSession.NewGameSession(
			(hostUser).(gameClient.IGameClient),
			createRoomReq.RoomName,
			createRoomReq.MinPlayer,
			createRoomReq.MaxPlayer,
			createRoomReq.WordList)
		//
		////Added host to player list of the session
		newSession.SetJoinedPlayer(hostUser.GetClientId(), hostUser)
		//
		//Added session into session manager
		gs.sessionManager.SetGameSession(newSession.GetSessionId(), newSession)
		//Sending a packet to client to updateUI?

		host := newSession.GetSessionHost()
		sessionName := newSession.GetSessionName()
		minPlayer, maxPlayer, currentPlayer := newSession.GetSessionPlayerInfo()
		status := newSession.GetSessionStatus()

		sessionInfo := packet.GameRoomInfoPacket{
			RoomId:            newSession.GetSessionId(),
			RoomName:          sessionName,
			RoomHostName:      host.GetName(),
			RoomHostId:        host.GetClientId(),
			RoomMinPlayer:     minPlayer,
			RoomMaxPlayer:     maxPlayer,
			RoomCurrentPlayer: currentPlayer,
			RoomStatus:        status,
		}

		resp := packet.CreateRoomResp{
			Code:               code.SUCCESS,
			GameRoomInfoPacket: sessionInfo,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized sessionInfo error : ", err)
			return
		}

		respPacket := packet.NewPacket(packetType.CREATE_ROOM, dataBytes)

		dataBytes, err = serializex.Marshal(&respPacket)
		if err != nil {
			log.Println("Serialized packet error : ", err)
			return
		}

		log.Println("sending response to client")
		if err := utils.SendMessage(host.GetConn(), dataBytes); err != nil {
			log.Println(err)
			return
		}
		break
	case packetType.JOIN_ROOM:
		//Join to an existing room
		var joinRoomReq packet.JoinRoomReq
		if err := serializex.Unmarshal(pkData, &joinRoomReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		//Find the room
		user, ok := gs.clientManager.GetGameClient(joinRoomReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", joinRoomReq.UserId)
			return

		}

		session, ok := gs.sessionManager.GetGameSession(joinRoomReq.RoomId)
		if !ok {
			log.Printf("Session %s not exist\n", joinRoomReq.RoomId)
			return
		}

		session.SetJoinedPlayer(user.GetClientId(), user)
		//Sending a packet to client to updateUI?
		host := session.GetSessionHost()
		sessionName := session.GetSessionName()
		minPlayer, maxPlayer, currentPlayer := session.GetSessionPlayerInfo()
		status := session.GetSessionStatus()

		sessionInfo := packet.GameRoomInfoPacket{
			RoomId:            session.GetSessionId(),
			RoomName:          sessionName,
			RoomHostName:      host.GetName(),
			RoomHostId:        host.GetClientId(),
			RoomMinPlayer:     minPlayer,
			RoomMaxPlayer:     maxPlayer,
			RoomCurrentPlayer: currentPlayer,
			RoomStatus:        status,
		}

		resp := packet.JoinRoomResp{
			Code:               code.SUCCESS,
			GameRoomInfoPacket: sessionInfo,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized sessionInfo error : ", err)
			return
		}

		respPacket := packet.NewPacket(packetType.JOIN_ROOM, dataBytes)

		dataBytes, err = serializex.Marshal(&respPacket)
		if err != nil {
			log.Println("Serialized packet error : ", err)
			return
		}

		if err := utils.SendMessage(user.GetConn(), dataBytes); err != nil {
			log.Println(err)
			return
		}

		break
	case packetType.EXIT_ROOM:
		//Exit current joined room
		var exitRoomReq packet.ExitRoomReq
		if err := serializex.Unmarshal(pkData, &exitRoomReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		user, ok := gs.clientManager.GetGameClient(exitRoomReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", exitRoomReq.UserId)
			return

		}

		session, ok := gs.sessionManager.GetGameSession(exitRoomReq.RoomId)
		if !ok {
			log.Printf("Session %s not exist\n", exitRoomReq.RoomId)
			return
		}

		session.SetExitedPlayer(user.GetClientId())

		//Sending a packet to client to updateUI?

		resp := packet.ExitRoomResp{
			Code: code.SUCCESS,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized sessionInfo error : ", err)
			return
		}

		respPacket := packet.NewPacket(packetType.EXIT_ROOM, dataBytes)

		dataBytes, err = serializex.Marshal(&respPacket)
		if err != nil {
			log.Println("Serialized packet error : ", err)
			return
		}

		if err := utils.SendMessage(user.GetConn(), dataBytes); err != nil {
			log.Println(err)
			return
		}

		break
	case packetType.GET_SESSION_INFO:
		//To get current room and game status info
		var getSessionInfoReq packet.GetSessionInfoReq
		if err := serializex.Unmarshal(pkData, &getSessionInfoReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		user, ok := gs.clientManager.GetGameClient(getSessionInfoReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", getSessionInfoReq.UserId)
			return

		}

		session, ok := gs.sessionManager.GetGameSession(getSessionInfoReq.RoomId)
		if !ok {
			log.Printf("Session %s not exist\n", getSessionInfoReq.RoomId)
			return
		}

		host := session.GetSessionHost()
		sessionName := session.GetSessionName()
		minPlayer, maxPlayer, currentPlayer := session.GetSessionPlayerInfo()
		status := session.GetSessionStatus()

		sessionInfo := packet.GameRoomInfoPacket{
			RoomId:            session.GetSessionId(),
			RoomName:          sessionName,
			RoomHostName:      host.GetName(),
			RoomHostId:        host.GetClientId(),
			RoomMinPlayer:     minPlayer,
			RoomMaxPlayer:     maxPlayer,
			RoomCurrentPlayer: currentPlayer,
			RoomStatus:        status,
		}

		resp := packet.GetSessionInfoResp{
			Code:               code.SUCCESS,
			GameRoomInfoPacket: sessionInfo,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized sessionInfo error : ", err)
			return
		}

		respPacket := packet.NewPacket(packetType.GET_SESSION_INFO, dataBytes)

		dataBytes, err = serializex.Marshal(&respPacket)
		if err != nil {
			log.Println("Serialized packet error : ", err)
			return
		}

		if err := utils.SendMessage(user.GetConn(), dataBytes); err != nil {
			log.Println(err)
			return
		}

		break
	case packetType.START_GAME:
		break
	case packetType.END_GAME:
		break
	case packetType.UPDATE_GAME_STATUS:
		break
		//
	}
}

func (gs *GameServer) Close() error {
	fmt.Println("Server is closing.")
	close(gs.closedSessionChan)
	close(gs.closedClientChan)
	close(gs.registerClientChan)
	return gs.listener.Close()
}
