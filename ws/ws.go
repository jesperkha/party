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
	// Server types
	MsgQuestion MessageType = "question" // Asking a new question now, waiting for you answer
	MsgSetup    MessageType = "setup"    // Prompt setup questions
	MsgResults  MessageType = "results"  // Show results from question answers

	// Client types
	MsgReady MessageType = "ready" // I have finished setup and am ready to play

	// Host types
	MsgNextQuestion MessageType = "next" // Move on to next question
)

type Player struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ServerMessage struct {
	Type MessageType `json:"type"`

	// Setup
	Prompts []string `json:"prompts"`
	Players []Player `json:"players"`

	// Question
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type ClientMessage struct {
	clientId uint        // set upon recv
	Type     MessageType `json:"type"`

	// Answer
	Choice uint `json:"choice"` // 1-4

	// Ready
	MostLikely     string `json:"mostLikely"`
	WouldYouRather string `json:"wouldYouRather"`
	TakeAShot      uint   `json:"takeAShot"` // id
	BlindAnswer    string `json:"blindAnswer"`
}

type Server struct {
	host    *Client
	clients map[uint]*Client
	notif   *notifier.Notifier

	connect    chan *Client
	disconnect chan *Client
	out        chan ServerMessage
	in         chan ClientMessage

	ids uint
}

func NewServer() *Server {
	return &Server{
		notif:      notifier.New(),
		clients:    make(map[uint]*Client),
		connect:    make(chan *Client, 100),
		disconnect: make(chan *Client, 100),
		out:        make(chan ServerMessage, 100),
		in:         make(chan ClientMessage, 100),
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

			s.BeginGame()

		case client := <-s.disconnect:
			delete(s.clients, client.ID)
			log.Printf("Client %d disconnected: %s", client.ID, client.Name)

		case msg := <-s.out:
			log.Println("sending broadcast: ", msg)
			for _, client := range s.clients {
				if err := client.Conn.WriteJSON(msg); err != nil {
					log.Println("Error sending message to client:", err)
				}
			}
			if err := s.host.Conn.WriteJSON(msg); err != nil {
				log.Println("Error sending message to host")
			}

		case msg := <-s.in:
			switch msg.Type {
			case MsgNextQuestion:
				s.out <- ServerMessage{
					Type:     MsgQuestion,
					Question: "What is your name?",
				}

			default:
				log.Println("unhandled client message")
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

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "must have name", http.StatusBadRequest)
			return
		}

		log.Println("Client connected")
		s.ConnectNew(conn, name)
	}
}

func (s *Server) ConnectHostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}

		s.host = &Client{
			ID:   s.nextId(),
			Name: "HOST",
			Conn: conn,
		}

		go s.readPump(s.host, s.notif)
		log.Println("Host connected")
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
				return
			}

			switch msgType {
			case websocket.TextMessage:
				log.Printf("Received message from client %d: %s", client.ID, msg)

				var message ClientMessage
				if err := json.Unmarshal(msg, &message); err != nil {
					log.Println("failed to parse json message")
					return
				}

				message.clientId = client.ID
				s.in <- message

			case websocket.CloseMessage:
				log.Printf("Client %d closed the connection", client.ID)
				return
			}
		}
	}
}

func (s *Server) BeginGame() {
	players := []Player{}
	for _, client := range s.clients {
		players = append(players, Player{
			ID:   client.ID,
			Name: client.Name,
		})
	}

	setup := ServerMessage{
		Type: MsgSetup,
		Prompts: []string{
			"foo",
			"bar",
		},
		Players: players,
	}

	s.out <- setup
}

func (s *Server) nextId() uint {
	s.ids++
	return s.ids
}
