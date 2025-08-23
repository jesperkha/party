package ws

type MessageType string

const (
	// Server types
	MsgQuestion MessageType = "question" // Asking a new question now, waiting for you answer
	MsgSetup    MessageType = "setup"    // Prompt setup questions
	MsgResults  MessageType = "results"  // Show results from question answers
	MsgJoined   MessageType = "joined"   // A new player has joined

	// Client types
	MsgReady  MessageType = "ready"  // I have finished setup and am ready to play
	MsgAnswer MessageType = "answer" // Here is my answer for this question

	// Host types
	MsgNextQuestion MessageType = "next"  // Move on to next question
	MsgBegin        MessageType = "begin" // All players have joined, begin game
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

	// Joined
	Player Player `json:"player"`
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
