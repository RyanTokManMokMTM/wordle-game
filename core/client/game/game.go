package game

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/client/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/internal/config"
)

func Start(c config.Config) {
	gameClient := client.NewClient(c)
	gameClient.Run()
}
