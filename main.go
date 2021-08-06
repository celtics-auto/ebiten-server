package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/server"
)

var addr = flag.String("addr", ":3000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	srv := server.NewServer()
	go srv.Run()

	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		server.UpgradeConn(w, r, srv)
	})

	log.Println(fmt.Sprintf("Starting server on %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
