package config

import conf "github.com/RyanTokManMokMTM/wordle-game/core/common/types"

type Config struct {
	conf.ServerConf `yaml:",inline"`
	Round           uint     `yaml:"round"`
	WordList        []string `yaml:"wordList"`
}
