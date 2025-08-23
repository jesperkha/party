package ws

import "github.com/jesperkha/pipoker/app"

type MessageType string

const (
	// Server types
	MsgQuestion MessageType = "question" // Asking a new question now, waiting for you answer
	MsgSetup    MessageType = "setup"    // Prompt setup questions
	MsgResults  MessageType = "results"  // Show results from question answers
	MsgJoined   MessageType = "joined"   // A new player has joined
	MsgFinish   MessageType = "finish"   // Game is finished, show results

	// Client types
	MsgReady  MessageType = "ready"  // I have finished setup and am ready to play
	MsgAnswer MessageType = "answer" // Here is my answer for this question

	// Host types
	MsgNextQuestion MessageType = "next"  // Move on to next question
	MsgBegin        MessageType = "begin" // All players have joined, begin game
	MsgTimer        MessageType = "timer" // Timer ran out
)

type ServerMessage struct {
	Type MessageType `json:"type"`

	// Setup
	Prompts []string     `json:"prompts"`
	Players []app.Player `json:"players"`

	// Question
	Question app.Question `json:"question"`

	// Joined
	Player app.Player `json:"player"`

	// Results
	Results []int `json:"results"`
}

type ClientMessage struct {
	clientId uint        // set upon recv
	Type     MessageType `json:"type"`

	// Answer
	Choice uint `json:"choice"` // 1-4

	// Ready
	Prompts app.Prompts `json:"prompts"`
}
