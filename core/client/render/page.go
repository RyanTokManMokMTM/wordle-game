package render

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/color"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/code"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/notificationType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/renderEvent"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/status"
	"github.com/brettski/go-termtables"
	"log"
	"os"
	"strconv"
	"strings"
)

const header string = ` _    _ _________________ _     _____   _____  ___ ___  ________ _ 
| |  | |  _  | ___ \  _  \ |   |  ___| |  __ \/ _ \|  \/  |  ___| |
| |  | | | | | |_/ / | | | |   | |__   | |  \/ /_\ \ .  . | |__ | |
| |/\| | | | |    /| | | | |   |  __|  | | __|  _  | |\/| |  __|| |
\  /\  | \_/ / |\ \| |/ /| |___| |___  | |_\ \ | | | |  | | |___|_|
 \/  \/ \___/\_| \_|___/ \_____|____/   \____|_| |_|_|  |_|____/(_)`

func isExit(input string, existInput string) bool {
	if strings.Compare(input, existInput) == 0 {
		return true
	}
	return false
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func writeStringToScreen(w *bufio.Writer, message string) {
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

func flushScreen(writer *bufio.Writer) {
	err := writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func headerInfo(w *bufio.Writer) {
	writeStringToScreen(w, "======================================================================\n")
	writeStringToScreen(w, fmt.Sprintln(color.Blue+header+color.Reset))
	writeStringToScreen(w, "======================================================================\n")
}

func roomInfo(w *bufio.Writer, info packet.GameRoomInfoPacket) {
	writeStringToScreen(w, "======================================================================\n")
	writeStringToScreen(w, fmt.Sprintf("Room ID : %s\n", info.RoomId))
	writeStringToScreen(w, fmt.Sprintf("Room Host Name : %s\n", info.RoomHostName))
	writeStringToScreen(w, fmt.Sprintf("Room Host Id  : %s\n", info.RoomHostId))
	writeStringToScreen(w, "----------------------------------------------------------------------\n")
	writeStringToScreen(w, fmt.Sprintf("Room name : %s\n", info.RoomName))
	writeStringToScreen(w, fmt.Sprintf("Room game staus : %s\n", info.RoomStatus))
	writeStringToScreen(w, "======================================================================\n")
}

func renderRoomTable(w *bufio.Writer, rooms []packet.GameRoomInfoPacket) {
	table := termtables.CreateTable()

	table.AddHeaders("Room Id", "Room Name", "Room Host Name", "Game status")
	for _, room := range rooms {
		table.AddRow(room.RoomId, room.RoomName, room.RoomHostName, room.RoomStatus)
	}
	writeStringToScreen(w, table.Render())
	writeStringToScreen(w, "=========================================================\n")
}

func roomInfoPage(w *bufio.Writer, c client.IClient, isHost bool, info packet.GameRoomInfoPacket, callback func(mode uint, roomId string), sendingMessage func(roomId, message string)) {
	headerInfo(w)
	roomInfo(w, info)

	var selectedMode uint
	if isHost {
		writeStringToScreen(w, "Welcome, you have created a new room.\n")
		writeStringToScreen(w, "You can start the game or enter /q to leave, and chatting with any message.\n")
		writeStringToScreen(w, "/s: Start the game\n")
		flushScreen(w)
	} else {
		writeStringToScreen(w, "Welcome, You joined the room. Please waiting room's host to start the game.\n")
		writeStringToScreen(w, "/q: Exit the room\n")
		flushScreen(w)
	}

	for {
		input, ok := <-c.GetInput()
		if !ok {
			log.Fatal("input channel closed")
		}

		input = strings.Trim(input, "\r\n")
		if strings.Compare(input, "-1") == 0 {
			return //just return
		}

		if b := isExit(input, "/q"); b {
			selectedMode = 0
			break //exit
		}

		if isHost && strings.Compare("/s", input) == 0 {
			selectedMode = 1 //start
			break
		}
		writeStringToScreen(w, color.Yellow+fmt.Sprintf("[You] : %s\n", input)+color.Reset)
		flushScreen(w)
		go sendingMessage(info.RoomId, input) //Sending message.
	}
	callback(selectedMode, info.RoomId) // sending game signal
}

func mainPage(c client.IClient) int {
	ClearScreen()
	w := bufio.NewWriter(os.Stdout)

	headerInfo(w)
	writeStringToScreen(w, fmt.Sprintln("Selected a mode or /q to exit:"))
	writeStringToScreen(w, fmt.Sprintln("1: Create a room."))
	writeStringToScreen(w, fmt.Sprintln("2: Join a room."))
	flushScreen(w)

	//Get input
	input, ok := <-c.GetInput()
	if !ok {
		log.Fatal("input channel closed.")
	}
	input = strings.Trim(input, "\r\n")
	b := isExit(input, "/q")
	if b {
		os.Exit(0)
	}

	mode, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}

	return mode
}

func setUpClientPage(c client.IClient, setNameCallBack func(name string)) {
	w := bufio.NewWriter(os.Stdout)
	headerInfo(w)
	writeStringToScreen(w, "Welcome to the multi-player wordle game!\n")
	flushScreen(w)
	var userName string
	for {
		writeStringToScreen(w, "Before starting the game, please enter your name: ")
		flushScreen(w)
		input, ok := <-c.GetInput()
		if !ok {
			log.Fatal("input channel closed.")
		}
		userName = strings.Trim(input, "\r\n")
		if len(userName) != 0 {
			break
		}
		writeStringToScreen(w, color.Red+"Name should not be empty.\n"+color.Reset)
	}

	setNameCallBack(userName)
}

func createRoomPage(c client.IClient, callback func(roomName string, wordList []string)) {
	ClearScreen()
	w := bufio.NewWriter(os.Stdout)

	var roomName string
	var wordList []string

	headerInfo(w)
	writeStringToScreen(w, "Enter the room information or enter /q to leave.\n")
	flushScreen(w)

	for {
		writeStringToScreen(w, "Enter a room name:")
		flushScreen(w)

		input, ok := <-c.GetInput()
		if !ok {
			log.Fatal("client input channel closed")
		}

		roomName = strings.Trim(input, "\r\n")
		b := isExit(input, "/q")
		if b {
			go func() {
				c.SetRenderEvent(renderEvent.HOME_PAGE, nil)
				c.SetRenderEventName(renderEvent.HOME_PAGE)
			}()
			return
		}

		if len(roomName) != 0 {
			break
		}

		writeStringToScreen(w, color.Red+"Room name should not be empty.\n"+color.Reset)
	}

	writeStringToScreen(w, "Enter the wordList(seperated by ,)(optional):")
	flushScreen(w)
	input, ok := <-c.GetInput()
	if !ok {
		log.Fatal("client input channel closed")
	}
	input = strings.Trim(input, "\r\n")
	words := strings.Split(input, ",")
	for _, w := range words {
		wordList = append(wordList, strings.TrimSpace(w))
	}

	b := isExit(input, "/q")
	if b {
		go func() {
			c.SetRenderEvent(renderEvent.HOME_PAGE, nil)
			c.SetRenderEventName(renderEvent.HOME_PAGE)
		}()
		return
	}
	callback(roomName, wordList)
}

func createRoomResultPage(c client.IClient, resp *packet.BasicResponseType, callback func(mode uint, roomId string), sendingMessage func(roomId, message string)) {
	ClearScreen()
	switch resp.Code {
	case code.SUCCESS:
		var result packet.CreateRoomResp
		if err := serializex.Unmarshal(resp.Data, &result); err != nil {
			log.Fatal(err)
			return
		}

		w := bufio.NewWriter(os.Stdout)
		roomInfoPage(w, c, true, result.GameRoomInfoPacket, callback, sendingMessage)
		flushScreen(w)
		break
	case code.REQUEST_FAILED:
		requestErrorHandler(c, resp.Message)
		return
	}
}

func joinRoomPage(c client.IClient, resp *packet.BasicResponseType, callback func(roomId string, isLeave bool)) {
	ClearScreen()
	switch resp.Code {
	case code.SUCCESS:
		var result packet.GetRoomListInfoResp
		if err := serializex.Unmarshal(resp.Data, &result); err != nil {
			log.Fatal(err)
			return
		}
		w := bufio.NewWriter(os.Stdout)
		renderRoomTable(w, result.Rooms)
		flushScreen(w)

		var roomId string
		isLeave := false
		for {
			writeStringToScreen(w, "Enter a room id you want to join or /q to leave\n")
			writeStringToScreen(w, "In: ")
			flushScreen(w)

			input, ok := <-c.GetInput()
			if !ok {
				log.Fatal("input channel closed")
			}
			input = strings.Trim(input, "\r\n")
			b := isExit(input, "/q")
			if b {
				isLeave = true
				break
			}

			//Check input is in the list
			canJoin := false
			for _, r := range result.Rooms {
				if strings.Compare(input, r.RoomId) == 0 && r.RoomStatus == status.ROOM_STAUS_WAITING {
					canJoin = true
					break
				}
			}
			if !canJoin {
				writeStringToScreen(w, color.Red+"You can not join the room due to room id or room status.\n"+color.Reset)
				continue
			}

			roomId = input
			break
		}
		callback(roomId, isLeave)
		break
	case code.REQUEST_FAILED:
		requestErrorHandler(c, resp.Message)
		return
	}

}

func joinRoomResultPage(c client.IClient, resp *packet.BasicResponseType, callback func(mode uint, roomId string), sendingMessage func(roomId, message string)) {
	ClearScreen()
	switch resp.Code {
	case code.SUCCESS:
		w := bufio.NewWriter(os.Stdout)

		var result packet.JoinRoomResp
		if err := serializex.Unmarshal(resp.Data, &result); err != nil {
			log.Fatal(err)
			return
		}

		roomInfoPage(w, c, false, result.GameRoomInfoPacket, callback, sendingMessage)
		flushScreen(w)
		break
	case code.REQUEST_FAILED:
		requestErrorHandler(c, resp.Message)
		return
	}
}

func gameStartingPage() {
	ClearScreen()
	w := bufio.NewWriter(os.Stdout)
	headerInfo(w)
	writeStringToScreen(w, "Game is started.")
	flushScreen(w)
}

func endingGamePage(c client.IClient, resp *packet.BasicResponseType, callback func(roomId string)) {
	ClearScreen()
	switch resp.Code {
	case code.SUCCESS:
		var result packet.EndingGameResp
		if err := serializex.Unmarshal(resp.Data, &result); err != nil {
			log.Fatal(err)
			return
		}
		w := bufio.NewWriter(os.Stdout)

		headerInfo(w)
		writeStringToScreen(w, result.OutputColorASNI+string(result.Message)+color.Reset)
		writeStringToScreen(w, "Enter any key to leave.\n")
		flushScreen(w)
		_, ok := <-c.GetInput()
		if !ok {
			log.Fatal("input channel closed")
		}
		callback(result.RoomId)
		break
	case code.REQUEST_FAILED:
		requestErrorHandler(c, resp.Message)
		return
	}

}

func gamingOutPut(c client.IClient, resp *packet.BasicResponseType) {
	switch resp.Code {
	case code.SUCCESS:
		var result packet.PlayingGameResp
		if err := serializex.Unmarshal(resp.Data, &result); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Print(result.OutputColorASNI + string(result.GameMessage) + color.Reset)
		if result.IsWritable {
			c.SetIsWritable()
		}

		break
	case code.REQUEST_FAILED:
		requestErrorHandler(c, resp.Message)
		return
	}
}

func notificationOutput(c client.IClient, resp *packet.BasicResponseType) {
	switch resp.Code {
	case code.SUCCESS:
		var result packet.NotifyPlayer
		if err := serializex.Unmarshal(resp.Data, &result); err != nil {
			log.Fatal(err)
			return
		}

		switch result.Type {
		case notificationType.ROOM_CHAT:
			fmt.Print(color.Yellow + string(result.Message) + color.Reset)
			break
		case notificationType.SYS:
			fmt.Print(color.Red + string(result.Message) + color.Reset)
			break
		default:
			fmt.Print(string(result.Message))
		}

		break
	case code.REQUEST_FAILED:
		requestErrorHandler(c, resp.Message)
		return
	}
}

func requestErrorHandler(c client.IClient, msg string) {
	w := bufio.NewWriter(os.Stdout)
	writeStringToScreen(w, color.Red+msg+color.Reset+"\n")
	writeStringToScreen(w, "Enter any key to leave.\n")
	flushScreen(w)

	_ = <-c.GetInput()
	c.SetRenderEvent(renderEvent.HOME_PAGE, nil)
	c.SetRenderEventName(renderEvent.HOME_PAGE)
	return
}
