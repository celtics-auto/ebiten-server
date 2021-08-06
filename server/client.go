package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type player struct {
	Position vector `json:"position"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type vector struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type message struct {
	Address string `json:"address"`
	Text    []byte `json:"text"`
}

type UpdateJson struct {
	Message *message `json:"message"`
	Player  *player  `json:"player"`
}

//Client precisa ter um ponteiro pro server
type Client struct {
	Address net.Addr
	Conn    *websocket.Conn
	server  *Server
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

/*
	TODO:
	Criar métodos de WritePump e ReadPump
	Client precisa ter o server como um atributo
	Client precisa ter o channel Send
*/

func (c *Client) Update(broadcast chan UpdateJson) {
	for {
		uJson := UpdateJson{}
		err := c.Conn.ReadJSON(&uJson)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("failed to read update json: %v", err)
			}
			break
		}
		broadcast <- uJson
	}

}

/*
Faz um upgrade da conexão para um websocket e instancia um objeto do client
Essa func tem que verificar se esse endereço está conectado
*/
func UpgradeConn(w http.ResponseWriter, r *http.Request, s *Server) {
	conn, err := upgrader.Upgrade(w, r, nil)
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
	client := NewClient(conn)

	log.Println(fmt.Sprintf("%s has connected.", client.Address))
	// TODO: Rodar duas goroutines chamando os métodos WritePump e ReadPump
}
