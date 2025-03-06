package main

import (
	"log"
	"os"

	"github.com/dmytro-vovk/shta/internal/boot"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configPath := "config.yml"
	if cfgPath := os.Getenv("CONFIG"); cfgPath != "" {
		configPath = cfgPath
	}

	c, err := boot.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Webserver().Serve(); err != nil {
		log.Fatal(err)
	}
}
