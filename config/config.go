package config

import (
	"github.com/joshy-joy/essay-word-counter/constants"
	"gopkg.in/yaml.v3"
	"os"
)

type Cgf struct {
	WebScrapper struct {
		Count int `yaml:"count"`
	} `yaml:"webScrapperJob"`
	Tokenizer struct {
		Count int `yaml:"count"`
	} `yaml:"tokenizerJob"`
	External struct {
		Timeout int64 `yaml:"timeoutInSeconds"`
	} `yaml:"external"`
	DefaultFilePath string `yaml:"defaultFilePath"`
	ResultLength    int    `yaml:"resultLength"`
	WordMinLength   int    `yaml:"wordMinLength"`
}

var config *Cgf

func InitConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var cfg Cgf
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}
	config = &cfg
	return nil
}

func Get() Cgf {
	return *config
}

func SetFilePath(path string) {
	if path != constants.Empty {
		config.DefaultFilePath = path
	}
}

func SetTopN(count int) {
	if count != 0 {
		config.ResultLength = count
	}
}
