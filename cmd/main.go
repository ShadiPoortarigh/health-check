package main

import (
	"flag"
	"fmt"
	"health-check/config"
	"health-check/pkg/postgres"
	"os"
)

var configPath = flag.String("config", "config.json", "configuration setup in json format")

func main() {
	flag.Parse()
	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		configPath = &v
	}

	c := config.MustReadConfig(*configPath)

	fmt.Println("Reading configuration file:", c)

	db, err := postgres.SetDB(c)
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("database Successfully connected: ", db)

}

// go run cmd/main.go --config sample-config.json
