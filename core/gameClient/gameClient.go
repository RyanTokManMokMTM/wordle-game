package main

import (
	"flag"
	"fmt"
	"github.com/RyanTokManMokMTM/wordle-game/core/common/conf"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameClient/client"
	"github.com/RyanTokManMokMTM/wordle-game/core/gameClient/internal/config"
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
	fmt.Println(c)
	gameClient := client.NewClient(c)
	gameClient.Run()

}
