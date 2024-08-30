package game

import (
	"bufio"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/internal/logic"
	"os"
	"strconv"
	"strings"
)

func Start(c client.IClient) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Select a mode : \n1:start a game\n2:exit")
		selectedMode, _ := reader.ReadString('\n')
		selectedMode = strings.Replace(selectedMode, "\n", "", -1)
		flag, err := strconv.Atoi(selectedMode)
		if err != nil {
			fmt.Println("Please input an number")
			continue
		}
		switch flag {
		case 1:
			logic.GameStart(c)
		case 2:
			fmt.Println("Thank you for playing")
			return
		default:
			fmt.Println("Mode not supported")
			continue
		}
	}

}
