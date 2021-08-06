package server

import (
	"log"

	"github.com/celtics-auto/ebiten-server/client"
	"github.com/celtics-auto/ebiten-server/config"
	"github.com/gorilla/websocket"
)

type Server struct {
	clients    map[*client.Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *UpdateJson
	cfg        *config.Server
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

/*
	O método Run vai ficar observando os 3 channels do server
*/
func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

/*
	TODO:
	Instanciar channels que faltam
	Remover o metodo SendMessages
*/
func New(clients client.ClientsMap, cfg *config.Server) *Server {
	broadcast := make(chan client.UpdateJson)

	return &Server{
		broadcast,
		clients,
		cfg,
	}
}
