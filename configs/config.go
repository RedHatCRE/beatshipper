package configs

import (
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v3"
)

const ConfigurationPath = "/etc/gz-beat-shipper/gz-beat-shipper-conf.yml"

type Configuration struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Path     string `yaml:"path"`
	Registry string `yaml:"registry"`
	Recheck  string `yaml:"recheck"`
}

func (c *Configuration) GetConfiguration() *Configuration {
	yamlFile, err := ioutil.ReadFile(ConfigurationPath)

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, c)

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
