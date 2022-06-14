package tournament

import (
	"fmt"
	"sort"

	"github.com/slshen/sb/pkg/game"
)

func GroupByTournament(games []*game.Game) (res []*Group) {
	// sort games by date
	sort.Slice(games, func(i, j int) bool {
		d1 := games[i].Date
		d2 := games[j].Date
		return d1.Before(d2)
	})
	res = []*Group{createTournamentGroup(games[0])}
	for _, g := range games[1:] {
		last := res[len(res)-1]
		if isSameTournament(last, g) {
			last.Games = append(last.Games, g)
			continue
		}
		res = append(res, createTournamentGroup(g))
	}
	return
}

func isSameTournament(gr *Group, g *game.Game) bool {
	d := gr.Games[len(gr.Games)-1].Date
	return g.Date == d || g.Date == d.AddDate(0, 0, 1)
}

func createTournamentGroup(g *game.Game) *Group {
	d := g.Date.Format("01/02/2006")
	var name string
	if g.Tournament != "" {
		name = fmt.Sprintf("%s %s", d, g.Tournament)
	} else {
		name = fmt.Sprintf("%s %s", d, g.League)
	}
	return &Group{
		Key:   g.Date,
		Name:  name,
		Games: []*game.Game{g},
	}
}
