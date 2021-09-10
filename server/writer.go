package server

import (
	"encoding/binary"

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

			zap.S().Debugf("x: %d - y: %d", message.Player.Position.X, message.Player.Position.Y)
			posX := uint16(message.Player.Position.X)
			posY := uint16(message.Player.Position.Y)
			buf := make([]byte, 4)
			binary.LittleEndian.PutUint16(buf[0:], posX)
			binary.LittleEndian.PutUint16(buf[2:], posY)

			if err := c.Conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					zap.S().Errorf("failed to write update json: %s", err.Error())

					c.server.unregister <- c
				}
			}
		}
	}
}
