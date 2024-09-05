package game

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/server/core/gameServer"
	"github.com/RyanTokManMokMTM/wordle-game/core/server/internal/config"
	"log"
)

func Start(c config.Config) {
	server := gameServer.NewGameServer(c)
	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}
