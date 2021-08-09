package server

import (
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		// TODO: implement a case where a ping is sent to the client every X time, if no response is sent within the deadline back kill it
		select {
		case message, ok := <-c.send:
			if !ok {
				// server closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Println("sending message")
			if err := c.Conn.WriteJSON(message); err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("failed to write update json:", err)

					c.server.unregister <- c
				}
			}
		}
	}
}
