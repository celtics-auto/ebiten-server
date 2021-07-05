package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/client"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var clients = client.NewMap()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func listenMessages(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func connection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	_, alreadyConn := clients.FindClient(conn.RemoteAddr())
	if alreadyConn {
		log.Println(fmt.Sprintf("Address %s already connected", conn.RemoteAddr()))
		return
	}
	cli := client.New(conn)
	defer clients.Disconnect(cli.Address)

	log.Println(fmt.Sprintf("%s has connected.", cli.Address))
	cli.ListenMessages()
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/connection", connection)

	log.Println(fmt.Sprintf("Starting server on %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
