package logic

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/regex"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/utils"
	"log"
	"math"
	"os"
	"strings"
)

/*
		FLOW:
		IN: HELLO, output: _____ (pick either one form the final list to check)
		- [HELLO, WORLD, QUITE, FANCY, FRESH, PANIC, CRAZY, BUGGY].
		- HELLO -> 00000 (1st)
		- WORLD -> _?_0_ (2st)
		- QUITE -> ____? (4th)
		- FANCY -> _____ (5th)
	   	- FRESH -> __?_? (3rd)
		- PANIC -> _____ (5th)
	    - CRAZY -> _____ (5th)
		- BUGGY -> _____ (5th)

		Pick all word with lower score , in this case will be
		Final list : [FANCY,PANIC,CRAZY,BUGGY]
		lower score is the output?

		INPUT: WORLD, output: _____(pick either one form the final list to check)
		- FANCY: _____ (2rd)
		- PANIC: _____ (2rd)
		- CRAZY: _?___ (1st)
		- BUGGY: _____ (2rd)
		Final list : [FANCY,PANIC,BUGGY]

		INPUT: FRESH ,output: _____(pick either one form the final list to check)
		- FANCY: 0____
		- PANIC: _____
		- BUGGY: _____
		Final list : [PANIC,BUGGY]

		INPUT: CRAZY ,output: ?_?__(pick either one form the final list to check)
		- PANIC: _?__?
		- BUGGY: ____0
		Final list : [PANIC]

		normal playing.
*/

// GameStart start a game for that client
func GameStart(c client.IClient) {
	fmt.Println("Start a Wordle Game.")
	fmt.Println("------------------------------------------------------------------------------------")

	candidateList := c.GetWordList()
	totalRound := c.GetTotalRound()
	var currentRound uint = 0

	/* TODO: Using a array to indicated there is/are alphabets included in the guessing word
	   	- Array size: [a-z][A-Z] , total 52
	    - Index calculation: according to the ascii code table, 'a'-'a'=0 ,so arr[0] ='a', arr[1]= 'b', arr[2] = 'c' ,a[26]='A', a[27]='B',etc
	*/
	candidateListSize := len(candidateList)
	candidateCounterList := make([][]uint, candidateListSize)

	//handling counter of all words
	for i := 0; i < candidateListSize; i++ {
		candidateCounterList[i] = make([]uint, 52)
		for _, w := range candidateList[i] {
			index := utils.LetterIndexOf(w)
			if index < 0 {
				log.Fatal("Guessing word must include alphabets only")
				return
			}
			candidateCounterList[i][index]++
		}
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

		var result string
		var isWin bool
		candidateList, candidateCounterList, result, isWin = hostCheatingGameChecking(guessText, candidateList, candidateCounterList)
		fmt.Printf("Round %d result: %s\n", currentRound+1, result)
		fmt.Println("------------------------------------------------------------------------------------")
		if isWin {
			fmt.Println("Configuration! You successfully guess the word!")
			fmt.Println("------------------------------------------------------------------------------------")
			return
		}
		currentRound += 1
	}

	fmt.Printf("Game is over, your are failed to guess to word in round : %d\n", totalRound)
	fmt.Println("------------------------------------------------------------------------------------")
}

func hostCheatingGameChecking(input string, candidateList []string, candidateListCounter [][]uint) (updatedCandidateList []string, updatedCandidateCounterList [][]uint, finalResult string, isWin bool) {
	candidateListSize := len(candidateList)
	updatedCandidateList = make([]string, 0)
	updatedCandidateCounterList = make([][]uint, 0)
	isWin = false

	//If the candidate list has only 1 word, by as normal game
	if candidateListSize == 0 {
		return
	}
	if candidateListSize == 1 {
		result, hit, _ := processGameResultAndScore(input, candidateList[0], candidateListCounter[0])
		updatedCandidateList = candidateList               //No need to update
		updatedCandidateCounterList = candidateListCounter //No need to update
		finalResult = result
		if hit == 5 { //Hit the word!
			isWin = true
		}
		return
	}

	tempCandidateList := make([]string, candidateListSize)
	copy(tempCandidateList, candidateList)

	minHit := uint(math.MaxInt)
	minPresent := uint(math.MaxInt)
	candidateScore := make([][]uint, candidateListSize) //To know all candidate score
	for i := 0; i < candidateListSize; i++ {
		candidateScore[i] = make([]uint, 2)
	}

	candidateResult := make([]string, candidateListSize) //To store all candidate result

	//Calculating the result
	for i := 0; i < candidateListSize; i++ {
		tempCounter := make([]uint, 52)
		copy(tempCounter, candidateListCounter[i])

		result, hit, present := processGameResultAndScore(input, candidateList[i], tempCounter)
		if minHit > hit {
			minHit = hit
			minPresent = uint(math.MaxInt) //minimum hit is updated, need to reset minPresent
		}

		if minHit == hit {
			//MARK: Only current hit met minHit be able to update present
			/*
				For example, minHit: 2, minPresent: 1
				local hit = 3, present = 0

				Because of 0 < 1,
				If minPresent is always updated if present greater than minPresent, it will cause an unexpected result.
				Final result become(minHit:2 , minPresent:0) ,which result is unexpected.
				Expected result : minHit:2 , minPresent:1
			*/
			if minPresent > present {
				minPresent = present
			}
		}

		candidateResult[i] = result
		candidateScore[i][0] = hit
		candidateScore[i][1] = present
	}

	//TODO: We now know all score of candidates and the minScore
	//TODO: We need to keep all candidates with minScore, otherwise drop it from the list
	for i := 0; i < candidateListSize; i++ {
		//If current candidate scope is larger than min, drop it
		if candidateScore[i][0] > minHit { //According to the rule `More Hit will have higher scores.`
			continue
		} else if candidateScore[i][0] == minHit {
			//According to the rule `If the number of Hit is the same, more Present will have higher score.`
			if candidateScore[i][1] > minPresent {
				continue
			}
		}

		//Other case will be lowest candidate
		updatedCandidateList = append(updatedCandidateList, candidateList[i]) //this list will have the same score
		finalResult = candidateResult[i]

		//Update counter list
		updatedCandidateCounterList = append(updatedCandidateCounterList, candidateListCounter[i])
	}

	return
}

func processGameResultAndScore(in, candidate string, wordCounter []uint) (result string, hit, present uint) {
	counter := make([]uint, len(wordCounter))
	copy(counter, wordCounter)

	size := len(candidate)
	temp := []byte("_____")
	hit = 0
	present = 0

	//TODO: handling the correct letter and correct spot and updating the counter list
	for i := 0; i < size; i++ {
		index := utils.LetterIndexOf(rune(in[i]))
		if counter[index] == 0 {
			continue
		}

		if in[i] == candidate[i] {
			counter[index] -= 1
			temp[i] = '0'
			hit += 1 //1 HIT
		}

	}

	//TODO: Handling the other case
	for i := 0; i < size; i++ {
		index := utils.LetterIndexOf(rune(in[i]))
		if counter[index] == 0 {
			continue
		}
		temp[i] = '?'
		present += 1 //1 PRESENT
	}

	result = string(temp)
	return
}
