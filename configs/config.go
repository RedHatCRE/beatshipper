package configs

import (
	"log"
	"os"
	"reflect"

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

	checkConfigurationFields(*c)

	return c
}

// Check if all the fields of the configuration struct have value
// Check if the field key of the configuration exist
func checkConfigurationFields(c Configuration) {
	structFields := reflect.ValueOf(c)

	for i := 0; i < structFields.NumField(); i++ {
		if structFields.Field(i).IsZero() {
			log.Fatalf("%s field of struct empty.", structFields.Type().Field(i).Name)
		}
	}
}
