package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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
		zap.S().Errorf("upgrade connection error: %s", err.Error())
		return
	}

	zap.S().Infof("%s connected", conn.RemoteAddr())
	client := NewClient(conn, s)
	s.register <- client

	go client.WritePump()
	go client.ReadPump()
}
