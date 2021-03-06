package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/config"
	"github.com/celtics-auto/ebiten-server/logger"
	"github.com/celtics-auto/ebiten-server/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err.Error())
	}

	lErr := logger.Init(&cfg.Logger, cfg.AppEnv)
	if lErr != nil {
		log.Fatalf("failed to initialize logger: %s", lErr.Error())
	}

	srv := server.New(&cfg.Server)
	go srv.Run()

	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		server.UpgradeConn(w, r, srv)
	})

	host := fmt.Sprintf("%s%s", ":", cfg.Server.Port)
	zap.S().Infof("Starting server on %s", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
