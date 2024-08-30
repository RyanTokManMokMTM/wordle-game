package config

import "github.com/RyanTokManMokMTM/wordle-game/core/common/types"

type Config struct {
	types.ServerConf `yaml:",inline"`
	Round            uint     `yaml:"round"`
	WordList         []string `yaml:"wordList"`
}
