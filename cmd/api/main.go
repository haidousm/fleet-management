package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/haidousm/fleets/internal/maps"
	"github.com/haidousm/fleets/internal/mqtt"
	"github.com/haidousm/fleets/internal/vcs"
)

type config struct {
	port  int
	env   string
	debug bool
}

type application struct {
	config config
	logger *slog.Logger
}

var (
	version = vcs.Version()
)

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	boundary_lines := []maps.Line{
		{Start: maps.Location{X: 1, Y: 1}, End: maps.Location{X: 499, Y: 1}},
		{Start: maps.Location{X: 1, Y: 1}, End: maps.Location{X: 1, Y: 499}},
		{Start: maps.Location{X: 499, Y: 1}, End: maps.Location{X: 499, Y: 499}},
		{Start: maps.Location{X: 1, Y: 499}, End: maps.Location{X: 499, Y: 499}},
	}

	hallway_lines := []maps.Line{
		{Start: maps.Location{X: 100, Y: 0}, End: maps.Location{X: 100, Y: 350}},
		{Start: maps.Location{X: 100, Y: 350}, End: maps.Location{X: 150, Y: 350}},
		{Start: maps.Location{X: 150, Y: 350}, End: maps.Location{X: 150, Y: 250}},
		{Start: maps.Location{X: 150, Y: 250}, End: maps.Location{X: 250, Y: 250}},
		{Start: maps.Location{X: 250, Y: 250}, End: maps.Location{X: 250, Y: 350}},
		{Start: maps.Location{X: 250, Y: 350}, End: maps.Location{X: 350, Y: 350}},
		{Start: maps.Location{X: 350, Y: 350}, End: maps.Location{X: 350, Y: 250}},
		{Start: maps.Location{X: 350, Y: 250}, End: maps.Location{X: 450, Y: 250}},
		{Start: maps.Location{X: 450, Y: 250}, End: maps.Location{X: 450, Y: 350}},
		{Start: maps.Location{X: 450, Y: 350}, End: maps.Location{X: 500, Y: 350}},
	}

	floor_map := maps.Map{
		Lines: append(boundary_lines, hallway_lines...),
		Size:  maps.Size{Width: 500, Height: 500},
	}

	go broadcastMap(floor_map)

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func broadcastMap(floor_map maps.Map) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	topic := "maps/floor"
	client := mqtt.Client(ctx, topic)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			jsonBytes, err := json.Marshal(floor_map)
			if err != nil {
				fmt.Printf("failed to marshal floor map: %s\n", err)
			}
			_, err = client.Publish(context.Background(), &paho.Publish{
				Topic:   topic,
				QoS:     0,
				Payload: jsonBytes,
			})
			if err != nil {
				fmt.Printf("failed to publish floor map: %s\n", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
