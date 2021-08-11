package server

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

			if err := c.Conn.WriteJSON(message); err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					zap.S().Errorf("failed to write update json: %s", err.Error())

					c.server.unregister <- c
				}
			}
		}
	}
}
