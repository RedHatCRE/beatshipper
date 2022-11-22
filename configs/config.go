package configs

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	Paths    []string `yaml:"path"`
	Registry string   `yaml:"registry"`
	Recheck  string   `yaml:"recheck"`
}

func (c *Configuration) GetConfiguration() *Configuration {
	userDirConfig, err := os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("beatshipper-conf")
	viper.AddConfigPath("/etc/beatshipper/")
	viper.AddConfigPath(userDirConfig)

	err = viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Configuration) GetRecheckInSeconds() int32 {
	t, err := time.ParseDuration(c.Recheck)

	if err != nil {
		log.Fatal(err)
	}

	s := t.Seconds()
	return int32(s)
}

func (c *Configuration) GetRecheckDuration() time.Duration {
	recheckSeconds := c.GetRecheckInSeconds()
	recheckDuration := time.Second * time.Duration(recheckSeconds)
	return recheckDuration
}
