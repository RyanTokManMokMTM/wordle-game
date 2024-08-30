package main

import (
	"flag"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/conf"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/internal/config"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameServer/server/gameServer"
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

	server := gameServer.NewGameServer(&c)
	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}

}
