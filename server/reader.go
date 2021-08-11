package server

import (
	"go.uber.org/zap"
)

func (c *Client) ReadPump() {
	defer func() {
		c.server.unregister <- c
		c.Conn.Close()
	}()

	for {
		uJson := &UpdateJson{}
		err := c.Conn.ReadJSON(uJson)
		if err != nil {
			zap.S().Errorf("failed to read update json: %s", err.Error())
			// if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("failed to read update json: %v", err)
			// }
			break
		}
		c.server.broadcast <- uJson
	}
}
