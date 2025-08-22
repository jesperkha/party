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

type MessageType string

const (
	MsgQuestion MessageType = "question" // Asking a new question now, waiting for you answer
)

type Message struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}

type Answer struct {
	clientId uint
	Choice   uint `json:"choice"` // 1-4
}

type Server struct {
	clients map[uint]*Client
	notif   *notifier.Notifier

	connect    chan *Client
	disconnect chan *Client
	broadcast  chan Message
	answer     chan Answer

	ids uint
}

func NewServer() *Server {
	return &Server{
		notif:      notifier.New(),
		clients:    make(map[uint]*Client),
		connect:    make(chan *Client, 100),
		disconnect: make(chan *Client, 100),
		broadcast:  make(chan Message, 100),
		answer:     make(chan Answer, 100),
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
			log.Println("sending broadcast: ", msg)
			for _, client := range s.clients {
				if err := client.Conn.WriteJSON(msg); err != nil {
					log.Println("Error sending message to client:", err)
				}
			}

		case ans := <-s.answer:
			log.Printf("got answer from client %d: option %d", ans.clientId, ans.Choice)

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

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "must have name", http.StatusBadRequest)
			return
		}

		log.Println("Client connected")
		s.ConnectNew(conn, name)
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
				var answer Answer
				if err := json.Unmarshal(msg, &answer); err != nil {
					log.Println("failed to parse json message")
					continue
				}
				answer.clientId = client.ID
				s.answer <- answer

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
