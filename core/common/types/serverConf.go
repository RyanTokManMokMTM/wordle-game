package types

type ServerConf struct {
	Host        string `yaml:"host"`
	Port        uint   `yaml:"port"`
	NetworkType string `yaml:"networkType"`
}
