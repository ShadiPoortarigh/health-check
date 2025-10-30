package main

import (
	"flag"
	"fmt"
	"health-check/app"
	"health-check/config"
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

	a := app.MustNewApp(c)

	fmt.Println("app created: ", a)

}

// go run cmd/main.go --config sample-config.json
