package server

import (
	"github.com/gorilla/websocket"
)

type player struct {
	Position vector
	Width    int
	Height   int
}

type vector struct {
	X int16
	Y int16
}
type message struct {
	Address string
	Text    []byte
}

type UpdateJson struct {
	Message *message
	Player  *player
}

//Client precisa ter um ponteiro pro server
type Client struct {
	Conn   *websocket.Conn
	server *Server
	send   chan *UpdateJson
}

func NewClient(conn *websocket.Conn, server *Server) *Client {
	send := make(chan *UpdateJson)

	return &Client{
		Conn:   conn,
		server: server,
		send:   send,
	}
}
