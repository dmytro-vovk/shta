package main

import (
	"log"
	"os"

	"github.com/dmytro-vovk/shta/internal/boot"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configPath := "config.yml"
	if envPath, ok := os.LookupEnv("CONFIG"); ok {
		configPath = envPath
	}

	b, err := boot.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	b.Verifier()

	if err := b.Webserver().Serve(); err != nil {
		log.Fatal(err)
	}
}
