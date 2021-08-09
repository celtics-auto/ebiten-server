package server

import (
	"log"
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
			log.Printf("failed to read update json: %v", err)
			// if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("failed to read update json: %v", err)
			// }
			break
		}
		log.Println(string(uJson.Message.Text))
		c.server.broadcast <- uJson
	}
}
