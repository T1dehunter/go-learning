package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"log"
	"sync"
)

import (
	"github.com/joho/godotenv"
)

type Config struct {
	MongoUrl string `env:"MONGO_URL,required"`
}

var config *Config
var once sync.Once

func NewConfig() *Config {
	once.Do(func() {
		err := godotenv.Load("server/config/config.env")
		if err != nil {
			log.Fatalf("Error on build config: unable to load config.env file: %e", err)
		}
		config = &Config{}
		err = env.Parse(config)
		if err != nil {
			log.Fatalf("Error on build config: unable to parse ennvironment variables: %e", err)
		}
		fmt.Printf("Environment: %s\n", config.MongoUrl)
	})
	return config
}
