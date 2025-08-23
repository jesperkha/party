package app

type App struct {
	roundNr     int
	totalRounds int // Calculcated from number of questions
	players     map[uint]Player
	ready       int      // number of players ready
	answered    int      // number of players that have answered the current question
	curQ        Question // Current question
	questions   []Question
}

type Player struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Points  int    `json:"points"`
	prompts Prompts
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

type Prompts struct {
	MostLikely     string `json:"mostLikely"`
	WouldYouRather string `json:"wouldYouRather"`
	TakeAShot      uint   `json:"takeAShot"` // id
	BlindAnswer    string `json:"blindAnswer"`
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

	a := &App{
		players: pmap,
	}

	return a
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
	if len(a.questions) == 0 {
		return Question{
			Text:    "",
			Options: []Option{},
		}
	}

	q := a.questions[0]
	a.questions = a.questions[1:]

	a.answered = 0
	a.curQ = q
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
func (a *App) PlayerReady(id uint, prompts Prompts) {
	a.ready += 1
	p := a.players[id]
	p.prompts = prompts
	a.players[id] = p
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
	return len(a.questions) == 0
}
