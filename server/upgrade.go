package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*
	TODO: Essa func tem que verificar se esse endereço está conectado
*/
func UpgradeConn(w http.ResponseWriter, r *http.Request, s *Server) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	log.Printf("%s connected", conn.RemoteAddr())
	client := NewClient(conn, s)
	s.register <- client

	go client.WritePump()
	go client.ReadPump()
}
