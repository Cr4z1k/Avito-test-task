package dbconf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Dbname string `yaml:"dbname"`
	User   string `yaml:"user"`
}

func GetConnectionString() string {
	configData, err := os.ReadFile("./internal/config/dbconf/conf.yaml")
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.User, os.Getenv("dbPassword"), config.Dbname, config.Host, config.Port)

	return connectionString
}
