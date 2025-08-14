package main

import (
	"net/http"
	"time"

	"github.com/raphico/go-device-telemetry-api/internal/config"
	"github.com/raphico/go-device-telemetry-api/internal/db"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
	transporthttp "github.com/raphico/go-device-telemetry-api/internal/transport/http"
)

func main() {
	log := logger.New("[telemetry-api] ")

	cfg := config.Load()

	dbpool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer dbpool.Close()

	err = db.Migrate(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := transporthttp.NewRouter(log)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.Logger,
	}

	log.Info("Server started on " + cfg.HTTPAddr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
