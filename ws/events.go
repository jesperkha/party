package ws

import (
	"log"
	"math/rand"

	"github.com/jesperkha/pipoker/app"
)

func (s *Server) onClientConnected(client *Client) {
	s.clients[client.ID] = client

	// Notify host of new player
	s.host.Conn.WriteJSON(ServerMessage{
		Type:   MsgJoined,
		Player: app.Player{ID: client.ID, Name: client.Name},
	})

	log.Printf("Client %d connected: %s", client.ID, client.Name)
}

func (s *Server) onClientDisconnected(client *Client) {
	delete(s.clients, client.ID)
	log.Printf("Client %d disconnected: %s", client.ID, client.Name)
}

func (s *Server) onServerBroadcast(msg ServerMessage) {
	log.Println("sending broadcast: ", msg)
	for _, client := range s.clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("Error sending message to client:", err)
		}
	}

	if err := s.host.Conn.WriteJSON(msg); err != nil {
		log.Println("Error sending message to host")
	}
}

func (s *Server) onClientMessage(msg ClientMessage) {
	switch msg.Type {
	case MsgReady:
		s.app.PlayerReady(msg.clientId)
		if s.app.Ready() {
			s.beginRound()
		}

	case MsgBegin:
		s.beginGame()

	case MsgAnswer:
		s.app.Answer(msg.clientId, msg.Choice)
		if s.app.Answered() {
			s.showResults()
		}

	case MsgTimer, MsgNextQuestion:
		s.nextQuestion()

	default:
		log.Println("unhandled client message")
	}
}

func (s *Server) beginGame() {
	players := []app.Player{}
	for _, client := range s.clients {
		players = append(players, app.Player{
			ID:   client.ID,
			Name: client.Name,
		})
	}

	s.app = app.New(players)

	setup := ServerMessage{
		Type: MsgSetup,
		Prompts: []string{
			"Hvem er mest sannsynlig til å...",
			randomBlindQuestion(),
		},
		Players: players,
	}

	s.out <- setup
}

func randomBlindQuestion() string {
	questions := []string{
		"Noe man sier til småbarn",
	}

	n := rand.Intn(len(questions))
	return questions[n]
}

func (s *Server) beginRound() {
	s.nextQuestion()
}

func (s *Server) showResults() {
	s.host.Conn.WriteJSON(ServerMessage{
		Type:    MsgResults,
		Results: s.app.RoundResults(),
	})
}

func (s *Server) nextQuestion() {
	if s.app.Done() {
		s.host.Conn.WriteJSON(ServerMessage{
			Type:    MsgFinish,
			Players: s.app.Podium(),
		})
		return
	}

	q := s.app.NextQuestion()
	s.out <- ServerMessage{
		Type:     MsgQuestion,
		Question: q,
	}
}
