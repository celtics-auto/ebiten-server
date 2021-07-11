package client

import (
	"net"

	"github.com/gorilla/websocket"
)

type Client struct {
	Address net.Addr
	Conn    *websocket.Conn
}

func New(conn *websocket.Conn) *Client {
	return &Client{
		Address: conn.RemoteAddr(),
		Conn:    conn,
	}
}
