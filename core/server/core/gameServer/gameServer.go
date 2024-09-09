package gameServer

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/color"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/code"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/notificationType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/status"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameClient"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gamePlayer"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameRoom"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/manager"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/internal/config"
	"log"
	"net"
	"strings"
	"sync"
)

type GameServer struct {
	sync.Mutex
	host          string
	port          uint
	networkType   string
	listener      net.Listener
	round         uint
	dictionaryMap map[string]byte
	wordList      []string
	roomManager   manager.IGameRoomManager
	clientManager manager.IGameClientManager

	closedClientChan chan []byte //which client is disconnected.
	closedRoomChan   chan []byte //which room is closed
	gameOverRoomChan chan []byte

	messageChan        chan packet.BasicPacket
	registerClientChan chan gameClient.IGameClient
}

func NewGameServer(c config.Config, dictMap map[string]byte) IGameServer {
	return &GameServer{
		host:          c.Host,
		port:          c.Port,
		networkType:   c.NetworkType,
		round:         c.Round,
		wordList:      c.WordList,
		dictionaryMap: dictMap,
		roomManager:   manager.NewGameRoomManager(),   //Managing all room which is created
		clientManager: manager.NewGameClientManager(), //Managing all client who is connected

		registerClientChan: make(chan gameClient.IGameClient), //A new client is connected
		messageChan:        make(chan packet.BasicPacket),     //Received a message from client
		closedClientChan:   make(chan []byte),                 //Received clientId that client is disconnected
		closedRoomChan:     make(chan []byte),                 //Received roomId that room need to be closed
		gameOverRoomChan:   make(chan []byte),                 //Received roomId that game is finished(all player)
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
			//MARK: To get closed signal from game client
			_ = <-newClient.GetClosedEvent()
			userId := newClient.GetClientId()

			gs.closedClientChan <- []byte(userId)
		}()

		go func() {
			///MARK: To get packet message data game client
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

// eventListener listening to different channel
func (gs *GameServer) eventListener() {
	for {
		select {
		//To register a client
		case client, ok := <-gs.registerClientChan:
			if !ok {
				fmt.Println("register-client channel is closed")
				return
			}

			fmt.Println("New client registered")
			gs.clientManager.SetGameClient(client.GetClientId(), client)

			fmt.Println("Current user : ", len(gs.clientManager.GetGameClientList()))
			break

		//To remove and disconnect a client
		case clientId, ok := <-gs.closedClientChan:
			if !ok {
				fmt.Println("closed-client channel is closed")
				return
			}
			fmt.Println("client disconnected  : ", string(clientId))

			gs.clientManager.RemoveGameClient(string(clientId))
			rooms := gs.roomManager.GetGameRoomList()

			for _, s := range rooms {
				s, ok := gs.roomManager.GetGameRoom(s.GetRoomId())
				if !ok {
					return
				}

				cid := string(clientId)
				if ok := s.RemovePlayer(cid); ok {
					for _, p := range s.GetAllPlayer() {
						if strings.Compare(cid, p.GetClient().GetClientId()) != 0 {
							s.NotifyPlayerWithMessage(p, fmt.Sprintf("[SYSTEM] Player %s is disconnected.\n", p.GetClient().GetName()))
						}
					}
				}

				if len(s.GetAllPlayer()) == 0 {
					go func() {
						gs.closedRoomChan <- []byte(s.GetRoomId())
					}()
				}
			}

			break

		//To remove a room
		case id, ok := <-gs.closedRoomChan:
			if !ok {
				fmt.Println("closed-room channel is closed")
				return
			}
			gs.roomManager.RemoveGameRoom(string(id))
			break

		//To handling packet message from client
		case message, ok := <-gs.messageChan:
			if !ok {
				fmt.Println("closed-room channel is closed")
				return
			}
			gs.handleMessage(message)
			break

		case roomId, ok := <-gs.gameOverRoomChan:
			if !ok {
				fmt.Println("closed-room channel is closed")
				return
			}
			rId := string(roomId)
			fmt.Println("room is finished the game", rId)

			s, ok := gs.roomManager.GetGameRoom(rId)
			if !ok {
				fmt.Println("room not found")
				return
			}

			var winner []string
			var score uint = 0
			players := s.GetAllPlayer()
			for _, p := range players {
				playerScore := p.GetScore()
				if playerScore >= score {
					if score == playerScore { //can be more than 1 people
						winner = append(winner, p.GetClient().GetName())
						continue
					} else {
						winner = winner[:0] //reset
						winner = append(winner, p.GetClient().GetName())
					}
					score = playerScore
				}
			}

			scoreMessage := fmt.Sprintf("The game is over, the largest score is %d.\n", score)
			winnerMessage := fmt.Sprintf("The winners is/are [%s]\n", strings.Join(winner, ", "))
			message := fmt.Sprintf("%s%s", scoreMessage, winnerMessage)

			s.SetRoomStatus(status.ROOM_STAUS_WAITING)
			endingGameResp := packet.EndingGameResp{
				OutputColorASNI: color.Yellow,
				RoomId:          rId,
				Message:         []byte(message),
			}

			dataBytes, err := serializex.Marshal(&endingGameResp)
			if err != nil {
				log.Println(err)
				return
			}
			for _, p := range players {
				p.GetClient().SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.END_GAME, dataBytes)
			}

		}
	}
}

// handleMessage handing the packet from client
func (gs *GameServer) handleMessage(pk packet.BasicPacket) {
	pkType := pk.PacketType
	pkData := pk.Data

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
		u, ok := gs.clientManager.GetGameClient(createRoomReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", createRoomReq.UserId)
			return
		}

		roomWordList := make([]string, 0)
		copy(roomWordList, gs.wordList)

		roomWordList = append(gs.wordList, createRoomReq.WordList...)

		player := gamePlayer.NewPlayer(u)
		newRoom := gameRoom.NewGameRoom(
			player,
			createRoomReq.RoomName,
			roomWordList,
			gs.round,
			gs.dictionaryMap)

		newRoom.AddPlayer(u.GetClientId(), player)
		gs.roomManager.SetGameRoom(newRoom.GetRoomId(), newRoom)
		host := newRoom.GetRoomHost()
		roomName := newRoom.GetRoomName()
		roomStatus := newRoom.GetRoomStatus()

		roomInfo := packet.GameRoomInfoPacket{
			RoomId:       newRoom.GetRoomId(),
			RoomName:     roomName,
			RoomHostName: host.GetClient().GetName(),
			RoomHostId:   host.GetClient().GetClientId(),
			RoomStatus:   roomStatus,
		}

		resp := packet.CreateRoomResp{
			GameRoomInfoPacket: roomInfo,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized error : ", err)
			return
		}
		u.SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.CREATE_ROOM, dataBytes)
		break
	case packetType.ROOM_LIST_INFO:
		log.Println("Received an message of create room")
		var getRoomListInfoReq packet.GetRoomListInfoReq
		if err := serializex.Unmarshal(pkData, &getRoomListInfoReq); err != nil {
			log.Println("serialized create room req error : ", err)
			return
		}

		//MARK: Get User Client by id
		u, ok := gs.clientManager.GetGameClient(getRoomListInfoReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", getRoomListInfoReq.UserId)
			return
		}

		rooms := make([]packet.GameRoomInfoPacket, 0)
		roomList := gs.roomManager.GetGameRoomList()
		for _, room := range roomList {
			roomInfo := packet.GameRoomInfoPacket{
				RoomId:       room.GetRoomId(),
				RoomName:     room.GetRoomName(),
				RoomHostName: room.GetRoomHost().GetClient().GetName(),
				RoomHostId:   room.GetRoomHost().GetClient().GetClientId(),
				RoomStatus:   room.GetRoomStatus(),
			}
			rooms = append(rooms, roomInfo)
		}

		resp := packet.GetRoomListInfoResp{
			Rooms: rooms,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized error : ", err)
			return
		}

		u.SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.ROOM_LIST_INFO, dataBytes)
		break
	case packetType.JOIN_ROOM:
		//Join to an existing room
		var joinRoomReq packet.JoinRoomReq
		if err := serializex.Unmarshal(pkData, &joinRoomReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		//Find the room
		u, ok := gs.clientManager.GetGameClient(joinRoomReq.UserId)
		if !ok {
			log.Printf("Client %s not exist\n", joinRoomReq.UserId)
			return

		}

		room, ok := gs.roomManager.GetGameRoom(joinRoomReq.RoomId)
		if !ok {
			log.Printf("room %s not exist\n", joinRoomReq.RoomId)
			return
		}

		if strings.Compare(room.GetRoomStatus(), status.ROOM_STAUS_PLAYING) == 0 {
			//Room is now playing, can not join.
			u.SendToClient(code.REQUEST_FAILED, "Game room started, you cannot join.", packetType.JOIN_ROOM, nil)
			return
		}

		newPlayer := gamePlayer.NewPlayer(u)
		room.AddPlayer(u.GetClientId(), newPlayer)
		//Sending a packet to client to updateUI?
		host := room.GetRoomHost()
		roomName := room.GetRoomName()
		gameStatus := room.GetRoomStatus()

		roomInfo := packet.GameRoomInfoPacket{
			RoomId:       room.GetRoomId(),
			RoomName:     roomName,
			RoomHostName: host.GetClient().GetName(),
			RoomHostId:   host.GetClient().GetClientId(),
			RoomStatus:   gameStatus,
		}

		resp := packet.JoinRoomResp{
			GameRoomInfoPacket: roomInfo,
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Println("Serialized error : ", err)
			return
		}
		u.SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.JOIN_ROOM, dataBytes)

		for _, p := range room.GetAllPlayer() {
			if strings.Compare(p.GetClient().GetClientId(), u.GetClientId()) != 0 {
				//MARK: Sending to the other player that is joined
				message := fmt.Sprintf("[SYSTEM] Player %s is joined.\n", u.GetName())
				room.NotifyPlayerWithMessage(p, message)
			}
		}
		break
	case packetType.EXIT_ROOM:
		//Exit current joined room
		var exitRoomReq packet.ExitRoomReq
		if err := serializex.Unmarshal(pkData, &exitRoomReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		u, ok := gs.clientManager.GetGameClient(exitRoomReq.UserId)
		if !ok {
			log.Printf("Client %s not exist or removed\n", exitRoomReq.UserId)
			return

		}

		room, ok := gs.roomManager.GetGameRoom(exitRoomReq.RoomId)
		if !ok {
			log.Printf("Room %s not exist or removed\n", exitRoomReq.RoomId)
			return
		}

		room.RemovePlayer(u.GetClientId())
		u.SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.EXIT_ROOM, nil)

		if len(room.GetAllPlayer()) == 0 {
			log.Println("Empty player")
			go func() {
				gs.closedRoomChan <- []byte(room.GetRoomId())
			}()
		} else {
			for _, p := range room.GetAllPlayer() {
				if strings.Compare(p.GetClient().GetClientId(), u.GetClientId()) != 0 {
					//MARK: Sending to the other player that is joined
					message := fmt.Sprintf("[SYSTEM] Player %s is left.\n", u.GetName())
					room.NotifyPlayerWithMessage(p, message)
				}
			}
		}
		break
	case packetType.START_GAME:
		var gameStartReq packet.GameStartReq
		if err := serializex.Unmarshal(pkData, &gameStartReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		s, ok := gs.roomManager.GetGameRoom(gameStartReq.RoomId)
		if !ok {
			log.Println("Room not found")
			return
		}

		s.SetGuessingWord()
		s.SetRoomStatus(status.ROOM_STAUS_PLAYING)

		players := s.GetAllPlayer()

		for _, p := range players {
			p.GetClient().SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.START_GAME, pk.Data)
			go s.StartGame(p) //All User start at same time
		}

		go func() {
			_, ok := <-s.GetTheGameIsOver() //Waiting the room to finish the game
			if !ok {
				return
			}

			gs.gameOverRoomChan <- []byte(s.GetRoomId())
		}()
		break
	case packetType.ROOM_CHAT_MESSAGE:
		var gameMsgReq packet.GameRoomChatMessage
		if err := serializex.Unmarshal(pkData, &gameMsgReq); err != nil {
			log.Println("serialized join room req error : ", err)
			return
		}

		sender, ok := gs.clientManager.GetGameClient(gameMsgReq.UserId)
		if !ok {
			log.Println("Client not connected")
			return
		}

		s, ok := gs.roomManager.GetGameRoom(gameMsgReq.RoomId)
		if !ok {
			log.Println("Room not found")
			return
		}

		isHost := strings.Compare(sender.GetClientId(), s.GetRoomHost().GetClient().GetClientId()) == 0
		var msg string
		if isHost {
			msg = fmt.Sprintf("[%s (Room Host)] : %s\n", sender.GetName(), gameMsgReq.Message)
		} else {
			msg = fmt.Sprintf("[%s] : %s\n", sender.GetName(), gameMsgReq.Message)
		}

		players := s.GetAllPlayer()
		resp := packet.NotifyPlayer{
			Type:    notificationType.ROOM_CHAT,
			Message: []byte(msg),
		}

		dataBytes, err := serializex.Marshal(&resp)
		if err != nil {
			log.Print("serialize NotifyPlayer error : ", err)
			return
		}

		for _, p := range players {
			if strings.Compare(sender.GetClientId(), p.GetClient().GetClientId()) == 0 {
				continue
			}
			p.GetClient().SendToClient(code.SUCCESS, code.CodeToMessage(code.SUCCESS), packetType.GAME_NOTIFICATION, dataBytes)
		}

		break
	}
}

func (gs *GameServer) Close() error {
	fmt.Println("Server is closing.")
	close(gs.closedClientChan)
	close(gs.closedRoomChan)
	close(gs.messageChan)
	close(gs.registerClientChan)
	close(gs.gameOverRoomChan)
	return gs.listener.Close()
}
