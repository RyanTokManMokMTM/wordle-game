package logic

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/regex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"log"
	"net"
	"strings"
)

// GameLogic start a game for that gameClient
func GameLogic(guessingWord string, totalRound uint, conn net.Conn) {
	wordCounter := make([]uint, 52)
	var currentRound uint = 0

	/* TODO: Using a array to indicated there is/are alphabets included in the guessing word
	   	- Array size: [a-z][A-Z] , total 52
	    - Index calculation: according to the ascii code table, 'a'-'a'=0 ,so arr[0] ='a', arr[1]= 'b', arr[2] = 'c' ,a[26]='A', a[27]='B',etc
	*/
	fmt.Println("Current word: ", guessingWord)
	for _, w := range guessingWord {
		index := utils.LetterIndexOf(w)
		if index < 0 {
			log.Fatal("Guessing word must include alphabets only")
			return
		}
		wordCounter[index]++
	}

	reader := bufio.NewReader(conn)

	for currentRound < totalRound {
		err := writeMessage(conn, packet.NewPacket(packetType.IN_GAME, false, fmt.Sprintf("Round %d \n", currentRound+1)))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		err = writeMessage(conn, packet.NewPacket(packetType.IN_GAME, true, fmt.Sprintf("Input your guessing word:")))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		data := make([]byte, 256)
		n, err := reader.Read(data[:])
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		guessText := strings.Replace(string(data[:n]), "\n", "", -1)
		guessText = strings.TrimSpace(string(data[:n]))

		if !regex.Regex(guessText, regex.FiveLetterWordMatcher) {
			err = writeMessage(conn, packet.NewPacket(packetType.IN_GAME, false, "Input must be a 5-letter word\n---------------------\n"))
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			continue
		}

		result, isWin := guessingWordChecking(guessText, guessingWord, wordCounter)
		if isWin {
			err = writeMessage(conn, packet.NewPacket(packetType.IN_GAME, false, "Configuration! You win the game!\n"))
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			break
		}
		err = writeMessage(conn, packet.NewPacket(packetType.IN_GAME, false, fmt.Sprintf("Round %d result: %s\n---------------------\n", currentRound+1, result)))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		//c.SetWordHistory(guessText)
		currentRound += 1
	}
	err := writeMessage(conn, packet.NewPacket(packetType.IN_GAME, false, fmt.Sprintf("Game is over, current round guessing word is %s\n---------------------\n", guessingWord)))
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
}

func guessingWordChecking(in, guessingWord string, wordCounter []uint) (result string, isWin bool) {
	//Due to golang Slice send as a slice addr passing into function
	//Copy the original wordCounter to a new slice
	counter := make([]uint, len(wordCounter))
	copy(counter, wordCounter)

	result = ""
	isWin = false
	for i := 0; i < len(guessingWord); i++ {
		index := utils.LetterIndexOf(rune(in[i]))
		if counter[index] == 0 {
			//letter not exist in guessing word
			result += "_"
		} else {
			if guessingWord[i] == in[i] {
				//Current letter is matched and
				counter[index] -= 1
				result += "0"
			} else {
				//Current letter not matched but in other spot of the guessing word.
				result += "?"
			}
		}
	}
	if strings.Compare(result, "00000") == 0 {
		isWin = true
	}
	return
}

func writeMessage(conn net.Conn, data packet.Packet) error {
	dataBytes, err := serializex.Marshal(data)
	if err != nil {
		return err
	}

	msgLen := uint32(len(dataBytes))
	err = binary.Write(conn, binary.BigEndian, msgLen)
	if err != nil {
		return err
	}

	_, err = conn.Write(dataBytes)
	if err != nil {
		return err
	}

	return nil
}
