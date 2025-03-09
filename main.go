package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vjurczenia/actorfreq/actorfreq"
	"github.com/vjurczenia/mblg/mblg"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func keepAwakeOnRender() {
	renderURL := os.Getenv("RENDER_URL")
	if renderURL != "" {
		for {
			http.Get(renderURL)
			time.Sleep(60 * time.Second)
		}
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			slog.Error("Error loading .env file", "error", err)
		}
	}

	actorfreq.SetUpDB()
	actorfreq.AddHandlers("/actorfreq/")

	mblg.AddHandlers("/mblg/")

	http.HandleFunc("/ping", pingHandler)

	go keepAwakeOnRender()

	port := "8080"
	slog.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
