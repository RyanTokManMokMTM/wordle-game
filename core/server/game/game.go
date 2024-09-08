package game

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameServer"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/internal/config"
	"log"
)

func Start(c config.Config, dictList []string) {
	//Convert list to a map
	dictMap := make(map[string]byte)
	for _, w := range dictList {
		dictMap[w] = 1
	}

	server := gameServer.NewGameServer(c, dictMap)
	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}
