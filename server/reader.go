package server

import (
	"bytes"
	"encoding/binary"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (c *Client) ReadPump() {
	defer func() {
		c.server.unregister <- c
		c.Conn.Close()
	}()

	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil || msgType != websocket.BinaryMessage {
			zap.S().Errorf("failed to read update message: %s", err.Error())
			// if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("failed to read update json: %v", err)
			// }
			break
		}

		v := vector{}
		buf := bytes.NewReader(msg)
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			zap.S().Errorf("failed to decode byte array: %v", err)
			continue
		}

		uJson := &UpdateJson{
			Player: &player{
				Position: v,
			},
		}
		zap.S().Debugf("x: %d - y: %d", uJson.Player.Position.X, uJson.Player.Position.Y)
		c.server.broadcast <- uJson
	}
}
