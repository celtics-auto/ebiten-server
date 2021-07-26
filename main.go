package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/client"
	"github.com/celtics-auto/ebiten-server/config"
	"github.com/celtics-auto/ebiten-server/server"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	cfg, _ := config.New()
	clients := client.NewMap()
	srv := server.New(clients, &cfg.Server)
	go srv.SendMessages()

	http.HandleFunc("/connection", srv.ConnectClient)

	log.Println(fmt.Sprintf("Starting server on %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
