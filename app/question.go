package app

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func RandomBlindQuestion() string {
	qs := []string{
		"Dette sier jeg til mamma før jeg legger meg...",
		"Noe man sier til en skuffa trener:",
		"Motto til en Sigma ulv",
		"Quote fra en Facebook mom",
	}
	r := rand.Intn(len(qs))
	return qs[r]
}

func randomBlindPairing() string {
	qs := []string{
		"Før Hitler døde sa han...",
		"Noe man sier etter sex",
		"Etter begravelsen tenkte jeg...",
		"Jeg dro opp buksa og sa",
		"I fetteren min sin 10års bursdag sang jeg...",
	}
	r := rand.Intn(len(qs))
	return qs[r]
}

func Prompt() string {
	return "Hvem er mest sannsynlig til å..."
}

func (a *App) MakeQuestions() {
	qs := []Question{}

	// Add all "who is?" questions
	for _, p := range a.players {
		options := []Option{}
		ps := a.fourRandomPlayers()
		for _, p := range ps {
			options = append(options, Option{
				Text:  p.Name,
				owner: p.ID,
			})
		}

		text := strings.ToLower(p.prompts.MostLikely)
		if text[len(text)-1] == '?' {
			text = text[:len(text)-1]
		}
		qs = append(qs, Question{
			Text:    fmt.Sprintf("%s %s?", Prompt(), text),
			Options: options,
		})
	}

	{
		// Add would you rather
		ps := a.Podium()
		if len(ps)%2 != 0 {
			ps = ps[1:]
		}
		for i := 0; i < len(ps); i += 2 {
			p1 := ps[i]
			p2 := ps[i+1]

			qs = append(qs, Question{
				Text: "Vil du heller...",
				Options: []Option{
					{Text: p1.prompts.WouldYouRather, owner: p2.ID},
					{Text: p2.prompts.WouldYouRather, owner: p1.ID},
				},
			})
		}
	}

	{
		ps := a.Podium()
		if len(ps) < 4 {
			log.Fatal("must have more than 4 players")
		}
		for i := 0; i+3 < len(ps); i += 4 {
			p1 := ps[i]
			p2 := ps[i+1]
			p3 := ps[i+2]
			p4 := ps[i+3]

			q := Question{
				Text: randomBlindPairing(),
				Options: []Option{
					{Text: p1.prompts.BlindAnswer, owner: p1.ID},
					{Text: p2.prompts.BlindAnswer, owner: p2.ID},
					{Text: p3.prompts.BlindAnswer, owner: p3.ID},
					{Text: p4.prompts.BlindAnswer, owner: p4.ID},
				},
			}

			qs = append(qs, q)
		}
	}

	rand.Shuffle(len(qs), func(i, j int) { qs[i], qs[j] = qs[j], qs[i] })
	a.questions = qs
}

func (a *App) fourRandomPlayers() []Player {
	m := a.players
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	selected := keys[:4]

	players := []Player{}
	for _, k := range selected {
		players = append(players, a.players[k])
	}

	return players
}
