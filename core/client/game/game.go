package game

import (
	"github.com/RyanTokManMokMTM/wordle-game/core/client/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/internal/config"
	"github.com/RyanTokManMokMTM/wordle-game/core/client/render"
)

func Start(c config.Config) {
	gameClient := client.NewClient(c)
	pageRender := render.NewRender(gameClient)
	pageRender.Run()
	gameClient.Run()
}
