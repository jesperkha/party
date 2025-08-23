package app

type App struct {
	roundNr     int
	totalRounds int // Calculcated from number of questions
	players     map[uint]Player
	ready       int      // number of players ready
	answered    int      // number of players that have answered the current question
	curQ        Question // Current question
	done        bool     // Is the game over?
}

type Player struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Points int    `json:"points"`

	// Their questions for this round
	MostLikely     string
	WouldYouRather string
	TakeAShot      uint
	BlindAnswer    string
}

type Question struct {
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

type Option struct {
	Text  string `json:"text"` // Prompt raw text
	Votes int    `json:"votes"`
	owner uint   // Id of player owning this prompt
}

type Result struct {
	Option string
	Votes  int
}

func New(players []Player) *App {
	pmap := make(map[uint]Player)
	for _, p := range players {
		pmap[p.ID] = p
	}

	return &App{
		players: pmap,
	}
}

func (a *App) Player(id uint) Player {
	return a.players[id]
}

func (a *App) Podium() []Player {
	players := []Player{}
	for _, p := range a.players {
		players = append(players, p)
	}
	return players
}

// Get a randomly chosen question from the remaining pool of prompts.
func (a *App) NextQuestion() Question {
	q := Question{
		Text: "Hva sier Ole Brum?",
		Options: []Option{
			{Text: "Bæ", owner: 2},
			{Text: "Mø"},
			{Text: "Hei"},
			{Text: "Hade"},
		},
	}

	a.answered = 0
	a.curQ = q
	a.done = true
	return q
}

// Report a players answer to the current question.
func (a *App) Answer(playerId, choice uint) {
	a.curQ.Options[choice].Votes++
	a.answered++

	if a.curQ.Options[choice].owner == playerId {
		p := a.players[playerId]
		p.Points += 1
		a.players[playerId] = p
	}
}

// Get results from this round
func (a *App) RoundResults() []int {
	res := []int{}
	for _, opt := range a.curQ.Options {
		res = append(res, opt.Votes)
	}

	return res
}

// Mark a player as ready
func (a *App) PlayerReady(id uint) {
	a.ready += 1
}

// Is everyone ready to play?
func (a *App) Ready() bool {
	return a.ready == len(a.players)
}

// Has everyone answered?
func (a *App) Answered() bool {
	return a.answered == len(a.players)
}

func (a *App) Done() bool {
	return a.done
}
