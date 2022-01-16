package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Bitmap        uint64   `yaml:"BitMap"`
	WorkDir       string   `yaml:"WorkDirectory"`
	InheritSource []string `yaml:"InheritSource"`
}

var Cfg = Config{}

func init() {
	yml, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	if err := yaml.Unmarshal(yml, &Cfg); err != nil {
		fmt.Println(err)
	}
}

func (config *Config) HitIndex(i uint64) bool {
	return config.Bitmap&(1<<i) != 0
}
