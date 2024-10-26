package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Cgf struct {
	WebScrapper struct {
		Count int `yaml:"count"`
	} `yaml:"webScrapperJob"`
	Tokenizer struct {
		Count int `yaml:"count"`
	} `yaml:"tokenizerJob"`
	External struct {
		Timeout int64 `yaml:"timeout"`
	} `yaml:"external"`
	DefaultFilePath string `yaml:"defaultFilePath"`
	ResultLeght     int    `yaml:"resultLeght"`
}

var config *Cgf

func InitConfig() error {
	f, err := os.Open("../externals/config.yml")
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
	config.DefaultFilePath = path
}

func SetTopN(count int) {
	config.ResultLeght = count
}
