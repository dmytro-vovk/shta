package boot

import (
	"context"
	_ "database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmytro-vovk/shta/internal/boot/config"
	"github.com/dmytro-vovk/shta/internal/counter"
	"github.com/dmytro-vovk/shta/internal/db"
	"github.com/dmytro-vovk/shta/internal/fetcher"
	"github.com/dmytro-vovk/shta/internal/filter"
	"github.com/dmytro-vovk/shta/internal/periodic"
	"github.com/dmytro-vovk/shta/internal/storage"
	"github.com/dmytro-vovk/shta/internal/web"
	"github.com/dmytro-vovk/shta/internal/web/handlers"
	_ "github.com/lib/pq"
)

type Boot struct {
	container
	configPath string
}

func New(configPath string) (*Boot, error) {
	b := &Boot{
		configPath: configPath,
	}

	if err := b.loadConfig(); err != nil {
		return nil, err
	}

	go b.arm()

	return b, nil
}

func (b *Boot) arm() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Got %v, shutting down...", <-quit)
	b.container.shutdown()
	log.Printf("Shutdown complete")

	os.Exit(0)
}

func (b *Boot) Config() *config.Config {
	return b.Get("Config").(*config.Config) //nolint:forcetypeassert
}

func (b *Boot) loadConfig() error {
	s, err := config.Load(b.configPath)
	if err != nil {
		return err
	}

	log.Printf("Config loaded from %s", b.configPath)

	b.Set("Config", s, nil)

	return nil
}

func (b *Boot) Filter() *filter.Filter {
	const id = "URL Filter"
	if s, ok := b.Get(id).(*filter.Filter); ok {
		return s
	}

	s := filter.New(b.Storage(), b.Fetcher())

	b.Set(id, s, nil)

	return s
}

func (b *Boot) Fetcher() *fetcher.Fetcher {
	const id = "URL Fetcher"
	if s, ok := b.Get(id).(*fetcher.Fetcher); ok {
		return s
	}

	s := fetcher.New(b.Config().Settings.ConcurrencyLimit)

	b.Set(id, s, nil)

	return s
}

func (b *Boot) Router() http.Handler {
	const id = "Web Router"
	if s, ok := b.Get(id).(http.Handler); ok {
		return s
	}

	m := http.NewServeMux()
	h := handlers.New(b.Storage(), b.Filter())

	m.HandleFunc("POST /", h.AddURL)
	m.HandleFunc("GET /", h.LatestURLs)

	s := http.Handler(m)

	b.Set(id, s, nil)

	return s
}

func (b *Boot) Webserver() *web.Server {
	const id = "Web Server"
	if s, ok := b.Get(id).(*web.Server); ok {
		return s
	}

	cfg := b.Config().WebServer

	s := web.New(cfg.Listen, b.Router())

	b.Set(id, s, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Stop(ctx); err != nil {
			log.Printf("Error shutting down web server: %s", err)
		}
	})

	log.Printf("Web server will listen on %s", cfg.Listen)

	return s
}

func (b *Boot) Storage() *storage.Storage {
	const id = "Storage"
	if s, ok := b.Get(id).(*storage.Storage); ok {
		return s
	}

	s := storage.New(b.Database(), b.Counter())

	b.Set(id, s, nil)

	return s
}

func (b *Boot) Counter() *counter.Counter {
	const id = "Counter"
	if s, ok := b.Get(id).(*counter.Counter); ok {
		return s
	}

	s := counter.New(b.Config().Settings.TopURLs)

	b.Set(id, s, nil)

	return s
}

func (b *Boot) Database() *db.DB {
	const id = "Database"
	if s, ok := b.Get(id).(*db.DB); ok {
		return s
	}

	var (
		s   *db.DB
		err error
	)

	for {
		if s, err = db.New(b.Config().Database.DSN()); err != nil {
			log.Printf("Error connecting to database: %s", err)

			time.Sleep(5 * time.Second)

			continue
		}

		break
	}

	b.Set(id, s, nil)

	return s
}

func (b *Boot) Verifier() *periodic.Verifier {
	const id = "URL Verifier"
	if s, ok := b.Get(id).(*periodic.Verifier); ok {
		return s
	}

	log.Printf("Will verify URLs every %s", b.Config().Settings.VerifyEvery)

	s := periodic.New(b.Counter(), b.Fetcher(), b.Config().Settings.VerifyEvery)

	b.Set(id, s, nil)

	return s
}
