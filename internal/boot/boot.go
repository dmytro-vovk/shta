package boot

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Boot struct {
	container
	configPath string
}

func New(configPath string) *Boot {
	b := Boot{
		configPath: configPath,
	}

	go b.arm()

	return &b
}

func (b *Boot) arm() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Got %v, shutting down...", <-quit)
	b.container.shutdown()
	log.Printf("Shutdown complete")

	os.Exit(0)
}
