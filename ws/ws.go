package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jesperkha/notifier"
)

type Client struct {
	ID     uint
	Name   string
	Avatar uint
	Conn   *websocket.Conn
}

type broadcast struct {
	sender  uint // id of sender
	message Message
}

type Server struct {
	clients map[uint]*Client
	notif   *notifier.Notifier

	connect    chan *Client
	disconnect chan *Client
	broadcast  chan broadcast

	ids uint
}

func NewServer() *Server {
	return &Server{
		notif:      notifier.New(),
		clients:    make(map[uint]*Client),
		connect:    make(chan *Client),
		disconnect: make(chan *Client),
		broadcast:  make(chan broadcast),
	}
}

func (s *Server) Run(notif *notifier.Notifier) {
	log.Println("WebSocket server is running")
	done, finish := notif.Register()

	for {
		select {
		case client := <-s.connect:
			s.clients[client.ID] = client
			log.Printf("Client %d connected: %s", client.ID, client.Name)

		case client := <-s.disconnect:
			delete(s.clients, client.ID)
			log.Printf("Client %d disconnected: %s", client.ID, client.Name)

		case msg := <-s.broadcast:
			for _, client := range s.clients {
				if client.ID == msg.sender {
					continue // do not send to sender
				}
				if err := client.Conn.WriteJSON(msg.message); err != nil {
					log.Println("Error sending message to client:", err)
				}
			}

		case <-done:
			log.Println("WebSocket server shutting down")
			s.notif.NotifyAndWait()
			finish()
			return
		}
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) ConnectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}

		log.Println("Client connected")
		s.ConnectNew(conn, r.URL.Query().Get("name"))
	}
}

// Connect new client to the server, returns the client instance with a unique ID.
func (s *Server) ConnectNew(conn *websocket.Conn, name string) (*Client, error) {
	client := &Client{
		ID:     s.nextId(),
		Conn:   conn,
		Name:   name,
		Avatar: 0, // TODO: give random avatar
	}

	s.connect <- client
	go s.readPump(client, s.notif)

	return client, nil
}

func (s *Server) readPump(client *Client, notif *notifier.Notifier) {
	done, finish := notif.Register()

	defer func() {
		client.Conn.Close()
		s.disconnect <- client
		finish()
	}()

	for {
		select {
		case <-done:
			log.Println("Client disconnected:", client.ID)
			return

		default:
			msgType, msg, err := client.Conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				continue
			}

			switch msgType {
			case websocket.TextMessage:
				log.Printf("Received message from client %d: %s", client.ID, msg)
				var msgParsed Message
				if err := json.Unmarshal(msg, &msgParsed); err != nil {
					log.Println("failed to parse json message")
					continue
				}

				s.broadcast <- broadcast{
					sender:  client.ID,
					message: msgParsed,
				}

			case websocket.CloseMessage:
				log.Printf("Client %d closed the connection", client.ID)
				return
			}
		}
	}
}

func (s *Server) nextId() uint {
	s.ids++
	return s.ids
}
