package logic

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/regex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"log"
	"os"
	"strings"
)

// GameStart start a game for that client
func GameStart(c client.IClient) {
	fmt.Println("Start a Wordle Game.")
	fmt.Println("---------------------")

	c.SetGuessingWord()
	guessingWord := c.GetGuessingWord()
	wordCounter := make([]uint, 52)
	totalRound := c.GetTotalRound()
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

	reader := bufio.NewReader(os.Stdin)
	for currentRound < totalRound {
		fmt.Printf("Round %d \n", currentRound+1)
		fmt.Printf("Input your guessing word : ")
		guessText, _ := reader.ReadString('\n')
		guessText = strings.Replace(guessText, "\n", "", -1)
		guessText = strings.TrimSpace(guessText)

		if !regex.Regex(guessText, regex.FiveLetterWordMatcher) {
			fmt.Println("Input must be a 5-letter word ")
			continue
		}

		result, isWin := guessingWordChecking(guessText, guessingWord, wordCounter)
		if isWin {
			fmt.Println("Configuration! You win the game!")
			break
		}
		fmt.Printf("Round %d result: %s\n", currentRound+1, result)
		fmt.Println("---------------------")
		c.SetWordHistory(guessText)
		currentRound += 1
	}

	fmt.Printf("Game is over, current round guessing word is %s\n", guessingWord)
	fmt.Printf("Your guess words: %s\n", strings.Join(c.GetWordHistory(), ","))
	fmt.Println("---------------------")
	c.Reset()
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
