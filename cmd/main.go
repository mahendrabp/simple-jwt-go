package main

import (
	"fmt"
	"log"
	"os"
	"simple-jwt-go/configs/persist/postgres"
	"simple-jwt-go/pkg/server"

	"simple-jwt-go/pkg/config"
)

func main() {
	fmt.Printf("PID: %d\n", os.Getpid())

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config"
	}

	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "local"
	}

	cfg, err := config.LoadConfig(configPath, configName)
	if err != nil {
		fmt.Printf("error config file: %v\n", err)
		os.Exit(1)
	}

	dbConfig := postgres.NewConfig(
		cfg.DB.Driver,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSL,
	)
	db, err := postgres.NewClient(dbConfig)
	if err != nil {
		log.Fatalf("no db connection: %v", err)
	}
	defer db.Close()

	s := server.New(cfg, db)
	if err := s.Run(); err != nil {
		log.Panicf("server is not running: %v", err)
	}
}
