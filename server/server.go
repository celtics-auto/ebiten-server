package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/client"
	"github.com/gorilla/websocket"
)

type MessageJson struct {
	Address string `json:"address"`
	Message []byte `json:"message"`
}

// TODO: break server logic into separate files e.g.: chat messages, player position update, update clients
type Server struct {
	clients   client.ClientsMap
	upgrader  websocket.Upgrader
	broadcast chan client.UpdateJson
}

func (s *Server) ConnectClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	_, alreadyConn := s.clients.FindClient(conn.RemoteAddr())
	if alreadyConn {
		log.Println(fmt.Sprintf("Address %s already connected", conn.RemoteAddr()))
		return
	}
	client := s.clients.CreateClient(conn)
	s.clients.Add(client)
	defer s.clients.Disconnect(client.Address)

	log.Println(fmt.Sprintf("%s has connected.", client.Address))
	client.Update(s.broadcast)
}

func (s *Server) SendMessages() {
	for {
		msg := <-s.broadcast
		for _, c := range s.clients {
			if err := c.Conn.WriteJSON(msg); err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("failed to write update json:", err)
					c.Conn.Close()
					s.clients.Disconnect(c.Address)
				}
				continue
			}
		}
	}
}

func New(clients client.ClientsMap) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	broadcast := make(chan client.UpdateJson)

	return &Server{
		clients,
		upgrader,
		broadcast,
	}
}
