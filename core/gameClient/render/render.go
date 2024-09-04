package render

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"log"
	"os"
	"strconv"
	"strings"
)

var header string = ` _    _ _________________ _     _____   _____  ___ ___  ________ _ 
| |  | |  _  | ___ \  _  \ |   |  ___| |  __ \/ _ \|  \/  |  ___| |
| |  | | | | | |_/ / | | | |   | |__   | |  \/ /_\ \ .  . | |__ | |
| |/\| | | | |    /| | | | |   |  __|  | | __|  _  | |\/| |  __|| |
\  /\  | \_/ / |\ \| |/ /| |___| |___  | |_\ \ | | | |  | | |___|_|
 \/  \/ \___/\_| \_|___/ \_____|____/   \____|_| |_|_|  |_|____/(_)`

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func readString(r *bufio.Reader) (string, error) {
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(input, "\r\n"), nil
}

func writeStringToScreen(w *bufio.Writer, message string) {
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}

func flushScreen(writer *bufio.Writer) {
	err := writer.Flush()
	if err != nil {
		log.Println(err)
	}
}

func Header(w *bufio.Writer) {
	writeStringToScreen(w, "=========================================================\n")
	writeStringToScreen(w, fmt.Sprintln(header))
	writeStringToScreen(w, "=========================================================\n")
}

func RoomInfo(w *bufio.Writer, info packet.GameRoomInfoPacket) {
	writeStringToScreen(w, "=========================================================")
	writeStringToScreen(w, fmt.Sprintf("Room ID : %s\n", info.RoomId))
	writeStringToScreen(w, fmt.Sprintf("Room Host Name : %s\n", info.RoomHostName))
	writeStringToScreen(w, fmt.Sprintf("Room Host Id  : %s\n", info.RoomHostId))
	writeStringToScreen(w, "---------------------------------------------------------\n")
	writeStringToScreen(w, fmt.Sprintf("Room name : %s\n", info.RoomName))
	writeStringToScreen(w, fmt.Sprintf("Minimun player : %d\n", info.RoomMinPlayer))
	writeStringToScreen(w, fmt.Sprintf("Maximun player : %d\n", info.RoomMaxPlayer))
	writeStringToScreen(w, fmt.Sprintf("Current player : %d\n", info.RoomCurrentPlayer))
	writeStringToScreen(w, "=========================================================\n")
}

func RoomInfoPage(w *bufio.Writer, isHost bool, info packet.GameRoomInfoPacket) {
	Header(w)
	RoomInfo(w, info)
	if isHost {
		writeStringToScreen(w, "You are host, you can start the game by input S")
	} else {
		writeStringToScreen(w, "Welcome, You joined the room. Please waiting room's host to start the game.")
	}
}

func SetUpClientPage(setNameCallBack func(name string)) {
	w := bufio.NewWriter(os.Stdout)
	r := bufio.NewReader(os.Stdin)
	_, _ = w.Write([]byte(header + "\n"))
	_, _ = w.Write([]byte("Welcome to the multiply player game!\n"))
	_, _ = w.Write([]byte("Before starting the game, please enter your name:\n"))
	_ = w.Flush()

	input, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	userName := strings.Trim(input, "\r\n")
	log.Println(userName)
	setNameCallBack(userName)
}

func MainPage(callback func(mode int)) {
	ClearScreen()
	w := bufio.NewWriter(os.Stdout)
	r := bufio.NewReader(os.Stdin)

	Header(w)
	writeStringToScreen(w, fmt.Sprintln("Selected a mode:"))
	writeStringToScreen(w, fmt.Sprintln("1: Create a room."))
	writeStringToScreen(w, fmt.Sprintln("2: Join a room."))
	flushScreen(w)

	input, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input = strings.Trim(input, "\r\n")

	mode, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}

	callback(mode)
}

func CreateRoomPage(callback func(roomName string, minPlayer, maxPlayer uint, wordList []string)) {
	ClearScreen()
	w := bufio.NewWriter(os.Stdout)
	r := bufio.NewReader(os.Stdin)

	var roomName string
	var minPlayer uint
	var maxPlayer uint
	var wordList []string

	Header(w)
	writeStringToScreen(w, "Enter the room information.\n")
	writeStringToScreen(w, "Enter the room name:")
	flushScreen(w)
	for {
		input, err := readString(r)
		if err != nil {
			log.Println(err)
			writeStringToScreen(w, "Enter the room name:")
			flushScreen(w)
			continue
		}

		if err == nil {
			roomName = input
			break
		}

	}

	writeStringToScreen(w, "Enter the room minimum player:")
	flushScreen(w)
	for {
		input, err := readString(r)
		if err != nil {
			log.Println(err)
			writeStringToScreen(w, "Enter the room minimum player:")
			flushScreen(w)
			continue
		}

		inputNumber, err := strconv.Atoi(input)
		if err != nil {
			log.Println(err)
			writeStringToScreen(w, "Enter the room minimum player:")
			flushScreen(w)
			continue
		}

		if err == nil {
			minPlayer = uint(inputNumber)
			break
		}
	}

	writeStringToScreen(w, "Enter the room maximum player:")
	flushScreen(w)
	for {
		input, err := readString(r)
		if err != nil {
			log.Println(err)
			writeStringToScreen(w, "Enter the room maximum player:")
			flushScreen(w)
			continue
		}

		inputNumber, err := strconv.Atoi(input)
		if err != nil {
			log.Println(err)
			writeStringToScreen(w, "Enter the room maximum player:")
			flushScreen(w)
			continue
		}

		if err == nil {
			maxPlayer = uint(inputNumber)
			break
		}

	}
	writeStringToScreen(w, "Enter the wordList(seperated by ,):")
	flushScreen(w)
	for {
		input, err := readString(r)
		if err != nil {
			log.Println(err)
			writeStringToScreen(w, "Enter the wordList(seperated by ,):")
			flushScreen(w)
			continue
		}

		words := strings.Split(input, ",")
		for _, w := range words {
			wordList = append(wordList, strings.TrimSpace(w))
		}

		if err == nil {
			break
		}
	}

	callback(roomName, minPlayer, maxPlayer, wordList)
}

func CreateRoomResultPage(data []byte, clientId string) {
	ClearScreen()
	var result packet.CreateRoomResp
	if err := serializex.Unmarshal(data, &result); err != nil {
		log.Println(err)
		return
	}

	isHost := result.RoomHostId == clientId

	w := bufio.NewWriter(os.Stdout)
	RoomInfoPage(w, isHost, result.GameRoomInfoPacket)
	flushScreen(w)
}

func JoinRoomPage(data []byte) {
	ClearScreen()
}

func JoinRoomResultPage(data []byte) {
	ClearScreen()

}

func ExistRoomResultPage(data []byte) {
	ClearScreen()

}
