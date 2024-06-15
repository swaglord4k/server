package model

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigType string

const (
	PostgresConfig ConfigType = "pg"
	MysqlConfig    ConfigType = "mysql"
)

type Conf struct {
	Host     string `yaml:"HOST"`
	User     string `yaml:"USER"`
	Password string `yaml:"PASSWORD"`
	Dbname   string `yaml:"DBNAME"`
	Port     int16  `yaml:"PORT"`
}

func (c *Conf) GetConf(configType ConfigType) *Conf {
	directory, _ := os.Getwd()
	yamlFile, err := os.ReadFile(directory + fmt.Sprintf("/config/conf.%s.yaml", configType))
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
