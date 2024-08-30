package main

import (
	"flag"
	"github.com/RyanTokManMokMTM/wordle-game/core/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/conf"
	"github.com/RyanTokManMokMTM/wordle-game/core/game"
	"github.com/RyanTokManMokMTM/wordle-game/core/internal/config"
	"log"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	//Load config data from yaml file
	if err := conf.Load(*configFile, &c); err != nil {
		log.Fatal(err)
	}

	player := client.NewClient(c)
	game.Start(player)

}
