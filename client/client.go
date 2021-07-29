package client

import (
	"log"
	"net"

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

func (c *Client) Update(broadcast chan UpdateJson) {
	for {
		uJson := UpdateJson{}
		err := c.Conn.ReadJSON(&uJson)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("failed to read update json: %v", err)
			}
			break
		}
		broadcast <- uJson
	}
}
