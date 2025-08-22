package app

type App struct {
	roundNr     int
	totalRounds int // Calculcated from number of questions
	players     []Player
}

type Player struct {
	ID     uint
	Name   string
	Points int

	// Their questions for this round
	MostLikely     string
	WouldYouRather string
	TakeAShot      uint
	BlindAnswer    string
}

type Question struct {
	Text    string
	Options []Option
}

type Option struct {
	Text  string // Prompt raw text
	Owner uint   // Id of player owning this prompt
}

type Result struct {
	Option string
	Votes  int
}

func New(players []Player) *App {
	return &App{
		players: players,
	}
}

// Get the top three players in order by highest points (highest first).
func (a *App) Podium() []Player {
	return nil
}

// Get a randomly chosen question from the remaining pool of prompts.
func (a *App) NextQuestion() Question {
	return Question{}
}

// Report a players answer to the current question.
func (a *App) Answer(playerId, choice uint) {
}

// Get results from this round
func (a *App) RoundResults() []Result {
	return nil
}
