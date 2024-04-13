package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"port"`
		Dbname string `yaml:"dbname"`
		User   string `yaml:"user"`
	} `yaml:"postgres"`
}

func GetConnectionString() string {
	configData, err := os.ReadFile("./internal/config/conf.yaml")
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Postgres.User, os.Getenv("DB_PASS"), config.Postgres.Dbname, config.Postgres.Host, config.Postgres.Port)

	return connectionString
}
