package server

import (
	"github.com/gorilla/websocket"
)

type player struct {
	Position vector `json:"position"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type vector struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type message struct {
	Address string `json:"address"`
	Text    []byte `json:"text"`
}

type UpdateJson struct {
	Message *message `json:"message"`
	Player  *player  `json:"player"`
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
