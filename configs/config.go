package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Host      string   `yaml:"host"`
	Port      string   `yaml:"port"`
	Paths     []string `yaml:"path"`
	Registry  string   `yaml:"registry"`
	LogSource string   `yaml:"logsource"`
}

func (c *Configuration) GetConfiguration() *Configuration {
	userDirConfig, err := os.UserConfigDir()

	if err != nil {
		log.Print(err)
	}

	if len(userDirConfig) > 0 {
		viper.AddConfigPath(userDirConfig)
	}

	viper.SetConfigName("beatshipper-conf")
	viper.AddConfigPath("/etc/beatshipper/")

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
