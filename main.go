package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/config"
	"github.com/celtics-auto/ebiten-server/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err.Error())
	}
	srv := server.New(&cfg.Server)
	go srv.Run()

	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		server.UpgradeConn(w, r, srv)
	})

	host := fmt.Sprintf("%s%s", ":", "3000")
	log.Println(fmt.Sprintf("Starting server on %s", host))
	log.Fatal(http.ListenAndServe(host, nil))
}
