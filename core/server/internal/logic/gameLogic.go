package logic

import (
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/color"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/regex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/serializex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packet"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/packetType"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/types/score"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gamePlayer"
	"log"
	"net"
	"strings"
)

var gameDesc = `
=============================================================================================================================
Desc: Get 6 chances to guess a 5-letter word.
Players have six attempts to guess a five-letter word, 
with feedback given for each guess in the form of coloured tiles indicating when letters match or occupy the correct position.
MISSING LETTER(Not including any letter in the guessing word) : represent '_'
PRESENT LETTER(Included letter but wrong spot) : represent '?'
HITTING LETTER(Included letter and correct sport) : represent '0'
------------------------------------------------------------------------------
For example, (WORLD), Input: HELLO, Output: _?___
------------------------------------------------------------------------------

In our game, it will calculate a score for you each round, the sooner you guess the word, the higher score you gain.
Your final score will be the largest score you get in each round.
Enjoy you game! Good luck.
=============================================================================================================================
`

var divider = `=============================================================================================================================`

// GameLogic start a game for that client
func GameLogic(guessingWord string, totalRound uint, player gamePlayer.IGamePlayer) {
	wordCounter := make([]uint, 52)
	var currentRound uint = 0
	var gameScore uint = 0

	/* TODO: Using a array to indicated there is/are alphabets included in the guessing word
	   	- Array size: [a-z][A-Z] , total 52
	    - Index calculation: according to the ascii code table, 'a'-'a'=0 ,so arr[0] ='a', arr[1]= 'b', arr[2] = 'c' ,a[26]='A', a[27]='B',etc
	*/
	fmt.Printf(color.Yellow+"[DEV] Current word: %s\n"+color.Reset, guessingWord)
	for _, w := range guessingWord {
		index := utils.LetterIndexOf(w)
		if index < 0 {
			log.Fatal("Guessing word must include alphabets only")
			return
		}
		wordCounter[index]++
	}
	conn := player.GetClient().GetConn()

	err := writeMessage(conn, false, gameDesc, color.Yellow, packetType.PLAYING_GAME)
	if err != nil {
		log.Println(err)
		_ = conn.Close()
		return
	}

	for currentRound < totalRound {
		roundMessage := fmt.Sprintf("Round %d \n", currentRound+1)
		inputMessage := fmt.Sprintf("Pleas input your guess with 5-letter\n")

		err := writeMessage(conn, true, fmt.Sprintf("%s%s", roundMessage, inputMessage), color.Reset, packetType.PLAYING_GAME)
		if err != nil {
			log.Println(err)
			_ = conn.Close()
			return
		}

		input := <-player.GetClient().GetGameGuessingInput()
		log.Println(input)
		guessText := strings.Replace(string(input), "\n", "", -1)
		guessText = strings.TrimSpace(string(input))

		if !regex.Regex(guessText, regex.FiveLetterWordMatcher) {
			err := writeMessage(conn, false, fmt.Sprintf("Input must be a 5-letter word\n%s\n", divider), color.Red, packetType.PLAYING_GAME)
			if err != nil {
				log.Println(err)
				_ = conn.Close()
				return
			}
			continue
		}

		result, isWin := guessingWordChecking(guessText, guessingWord, wordCounter)
		currentScore := calculateUserScore(totalRound, currentRound, result)
		gameScore = max(currentScore, gameScore)
		player.SetScore(gameScore)
		if isWin {
			winMessage := fmt.Sprintf("Configuration! You guessed the word %s", guessingWord)
			scoreMessage := fmt.Sprintf("You final score is %d.\n", gameScore)

			err := writeMessage(conn, false, fmt.Sprintf("%s%s\n", winMessage, scoreMessage), color.Yellow, packetType.PLAYING_GAME)
			if err != nil {
				log.Println(err)
				_ = conn.Close()
				return
			}
			return
		}

		resultMessage := fmt.Sprintf("Round %d , result: %s\n", currentRound+1, result)
		currentScoreMessage := fmt.Sprintf("Round %d score is %d", currentRound+1, currentScore)

		err = writeMessage(conn, false, fmt.Sprintf("%s%s\n%s\n", resultMessage, currentScoreMessage, divider), color.Reset, packetType.PLAYING_GAME)
		if err != nil {
			log.Println(err)
			_ = conn.Close()
			return
		}

		currentRound += 1
	}

	gameOverMessage := fmt.Sprintf("Game is over, current round guessing word is %s\n", guessingWord)
	err = writeMessage(conn, false, fmt.Sprintf("%s\n", gameOverMessage), color.Red, packetType.PLAYING_GAME)
	if err != nil {
		log.Println(err)
		_ = conn.Close()
		return
	}

	player.SetScore(gameScore)
}

func calculateUserScore(totalRound, currenRound uint, currentResult string) uint {
	hit := 0
	miss := 0
	present := 0

	for _, c := range currentResult {
		switch c {
		case '_':
			miss += 1
			break
		case '?':
			present += 1
			break
		case '0':
			hit += 1
			break
		}
	}
	log.Println(hit, miss, present)

	return uint(int(totalRound-currenRound) * (hit*score.HIT + present*score.PRESENT + miss*score.MISS))
}

func guessingWordChecking(in, guessingWord string, wordCounter []uint) (result string, isWin bool) {
	//Due to golang Slice send as a slice addr passing into function
	//Copy the original wordCounter to a new slice
	counter := make([]uint, len(wordCounter))
	copy(counter, wordCounter)

	size := len(guessingWord)
	temp := []byte("_____")
	isWin = false

	/*
		Explain:
		If guessing word is "WORLD", user input is "HELLO"
		According to the algorithm which is using a counter.
		"WORLD" has only 1 letter 'L', "HELLO" have 2 letter 'L'
		The expected answers will be '___0?'

		If we're using 1 loop to handle this case, the answer is unexpected, which is '__?_?'

		To solve this issue:
		1. handling the correct letter and correct spot cases ,and updating the counter list
		2. handling the other case
	*/
	//TODO: handling the correct letter and correct spot and updating the counter list
	for i := 0; i < size; i++ {
		index := utils.LetterIndexOf(rune(in[i]))
		if counter[index] == 0 {
			continue
		}

		if in[i] == guessingWord[i] {
			counter[index] -= 1
			temp[i] = '0'
		}

	}

	//TODO: Handling the other case
	for i := 0; i < len(guessingWord); i++ {
		index := utils.LetterIndexOf(rune(in[i]))
		if counter[index] == 0 {
			continue
		}
		temp[i] = '?'
	}

	result = string(temp)
	if strings.Compare(result, "00000") == 0 {
		isWin = true
	}
	return
}

func writeMessage(conn net.Conn, isWritable bool, message string, colorASNI string, pkgType string) error {
	playingGameReq := packet.PlayingGameResp{
		OutputColorASNI: colorASNI,
		IsWritable:      isWritable,
		GameMessage:     []byte(message),
	}

	dataBytes, err := serializex.Marshal(&playingGameReq)
	if err != nil {
		return err
	}

	pk := packet.NewPacket(pkgType, dataBytes)
	dataBytes, err = serializex.Marshal(&pk)
	if err != nil {
		return err
	}

	return utils.SendMessage(conn, dataBytes)
}
