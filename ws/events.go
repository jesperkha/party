package ws

import "log"

func (s *Server) onClientConnected(client *Client) {
	s.clients[client.ID] = client

	// Notify host of new player
	s.host.Conn.WriteJSON(ServerMessage{
		Type:   MsgJoined,
		Player: Player{ID: client.ID, Name: client.Name},
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
	case MsgNextQuestion:
		s.out <- ServerMessage{
			Type:     MsgQuestion,
			Question: "What is your name?",
		}

	case MsgBegin:
		log.Println("starting game")

	default:
		log.Println("unhandled client message")
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
