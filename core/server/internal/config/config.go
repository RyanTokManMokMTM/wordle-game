package config

import "github.com/RyanTokManMokMTM/wordle-game/core/common/types/config"

type Config struct {
	config.ServerConf `yaml:",inline"`
	Round             uint     `yaml:"round"`
	WordList          []string `yaml:"wordList"`
}
