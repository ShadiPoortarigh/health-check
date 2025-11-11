package main

import (
	"flag"
	"fmt"
	"health-check/api/handlers/http"
	"health-check/app"
	"health-check/config"
	"log"
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

	log.Fatal(http.Run(a, c.Server))

}

// go run cmd/main.go --config sample-config.json
