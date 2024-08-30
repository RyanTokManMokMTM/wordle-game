package logic

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/regex"
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
		_, err := conn.Write([]byte(fmt.Sprintf("Round %d \n", currentRound+1)))
		if err != nil {
			conn.Close()
			return
		}
		_, err = conn.Write([]byte(fmt.Sprintf("Input your guessing word : ")))
		if err != nil {
			conn.Close()
			return
		}
		data := make([]byte, 256)
		n, err := reader.Read(data[:])
		if err != nil {
			conn.Close()
			return
		}
		guessText := strings.Replace(string(data[:n]), "\n", "", -1)
		guessText = strings.TrimSpace(string(data[:n]))

		if !regex.Regex(guessText, regex.FiveLetterWordMatcher) {
			_, err := conn.Write([]byte("Input must be a 5-letter word\n---------------------\n"))
			if err != nil {
				conn.Close()
				return
			}
			continue
		}

		result, isWin := guessingWordChecking(guessText, guessingWord, wordCounter)
		if isWin {
			_, err := conn.Write([]byte("Configuration! You win the game!\n"))
			if err != nil {
				conn.Close()
				return
			}
			break
		}
		msg := fmt.Sprintf("Round %d result: %s\n---------------------\n", currentRound+1, result)
		_, err = conn.Write([]byte(msg))
		if err != nil {
			conn.Close()
			return
		}

		//c.SetWordHistory(guessText)
		currentRound += 1
	}
	msg := []byte(fmt.Sprintf("Game is over, current round guessing word is %s\n---------------------", guessingWord))
	_, err := conn.Write(msg)
	if err != nil {
		conn.Close()
		return
	}
}

func guessingWordChecking(in, guessingWord string, wordCounter []uint) (result string, isWin bool) {
	result = ""
	isWin = false
	for i := 0; i < len(guessingWord); i++ {
		index := utils.LetterIndexOf(rune(in[i]))
		if wordCounter[index] == 0 {
			//letter not exist in guessing word
			result += "_"
		} else {
			if guessingWord[i] == in[i] {
				//Current letter is matched and
				wordCounter[index] -= 1
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
