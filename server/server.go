package server

import (
	"github.com/celtics-auto/ebiten-server/config"
)

type Server struct {
	broadcast  chan *UpdateJson
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
	cfg        *config.Server
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

func New(cfg *config.Server) *Server {
	broadcast := make(chan *UpdateJson)
	register := make(chan *Client)
	unregister := make(chan *Client)
	clients := make(map[*Client]bool)

	return &Server{
		broadcast,
		register,
		unregister,
		clients,
		cfg,
	}
}
