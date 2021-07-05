package client

import (
	"log"
	"net"

	"github.com/gorilla/websocket"
)

type Client struct {
	Address net.Addr
	conn    *websocket.Conn
}

func New(conn *websocket.Conn) *Client {
	return &Client{
		Address: conn.RemoteAddr(),
		conn:    conn,
	}
}

func (c *Client) ListenMessages() {
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
