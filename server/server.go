package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/celtics-auto/ebiten-server/client"
	"github.com/gorilla/websocket"
)

type Message struct {
	mType int
	data  []byte
}

type Server struct {
	clients  client.ClientsMap
	upgrader websocket.Upgrader
	msgChan  chan Message
}

func (s *Server) ConnectClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	_, alreadyConn := s.clients.FindClient(conn.RemoteAddr())
	if alreadyConn {
		log.Println(fmt.Sprintf("Address %s already connected", conn.RemoteAddr()))
		return
	}
	client := s.clients.CreateClient(conn)
	s.clients.Add(client)
	defer s.clients.Disconnect(client.Address)

	log.Println(fmt.Sprintf("%s has connected.", client.Address))
	s.ListenMessages(client)
}

func (s *Server) ListenMessages(c *client.Client) {
	for {
		mt, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("%s sent: %s", c.Address, message)

		msg := Message{
			mType: mt,
			data:  message,
		}
		s.msgChan <- msg
	}
}

func (s *Server) SendMessages() {
	for {
		msg := <-s.msgChan

		for _, c := range s.clients {
			fmt.Println(msg.data)
			if err := c.Conn.WriteMessage(msg.mType, msg.data); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func New(clients client.ClientsMap) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	msgChan := make(chan Message)

	return &Server{
		clients,
		upgrader,
		msgChan,
	}
}
