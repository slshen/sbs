package stats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/slshen/sb/pkg/game"
)

type GameStats struct {
	TeamStats map[string]*Stats
	OnlyTeam  string
}

func NewGameStats() *GameStats {
	return &GameStats{
		TeamStats: make(map[string]*Stats),
	}
}

func (mg *GameStats) Read(g *game.Game) error {
	states := g.GetStates()
	for i, state := range states {
		var battingTeam, fieldingTeam *game.Team
		if state.Top() {
			battingTeam = g.VisitorTeam
			fieldingTeam = g.HomeTeam
		} else {
			battingTeam = g.HomeTeam
			fieldingTeam = g.VisitorTeam
		}
		battingTeamStats := mg.GetStats(battingTeam)
		fieldingTeamStats := mg.GetStats(fieldingTeam)
		battingTeamStats.RecordBatting(g, state)
		fieldingTeamStats.RecordPitching(g, state, lastState(states, i))
		for _, runnerID := range state.ScoringRunners {
			runner := battingTeamStats.GetBatting(runnerID)
			runner.RunsScored++
		}
		switch state.Play.Type {
		case game.StolenBase:
			for _, runnerID := range state.Play.Runners {
				runner := battingTeamStats.GetBatting(runnerID)
				runner.StolenBases++
			}
		case game.CaughtStealing:
			if !state.NotOutOnPlay {
				runner := battingTeamStats.GetBatting(state.Play.Runners[0])
				runner.CaughtStealing++
			}
		}
	}
	return nil
}

func lastState(states []*game.State, i int) *game.State {
	if i == 0 {
		return nil
	}
	state := states[i]
	lastState := states[i-1]
	if state.InningNumber == lastState.InningNumber && state.Half == lastState.Half {
		return lastState
	}
	return nil
}

func (mg *GameStats) GetStats(team *game.Team) *Stats {
	stats := mg.TeamStats[team.Name]
	if stats == nil {
		stats = NewStats(team)
		mg.TeamStats[team.Name] = stats
	}
	return stats
}

type dataMaker struct {
	columnIndexes map[string]int
	data          Data
}

func (dm *dataMaker) addRow(m map[string]interface{}) {
	if dm.columnIndexes == nil {
		dm.data.Width = map[string]int{
			"Name": 20,
		}
		dm.columnIndexes = map[string]int{
			"Name": 0,
		}
		dm.data.Columns = []string{"Name"}
		for _, k := range sortKeys(m) {
			if _, ok := dm.columnIndexes[k]; !ok {
				dm.columnIndexes[k] = len(dm.columnIndexes)
				dm.data.Columns = append(dm.data.Columns, k)
			}
		}
	}
	row := make([]interface{}, len(dm.columnIndexes))
	for k, v := range m {
		row[dm.columnIndexes[k]] = v
	}
	dm.data.Rows = append(dm.data.Rows, row)
}

func (mg *GameStats) GetPitchingData() *Data {
	var dm dataMaker
	for team, stats := range mg.TeamStats {
		if mg.filterExclude(team, stats) {
			continue
		}
		var players []game.PlayerID
		for player := range stats.Pitching {
			players = append(players, player)
		}
		for _, player := range sortPlayers(players) {
			pitching := stats.Pitching[player]
			var m map[string]interface{}
			if err := mapstructure.Decode(pitching, &m); err != nil {
				panic(err)
			}
			mg.adjustRowValues(len(pitching.Games), team, pitching.Player, m)
			dm.addRow(m)
		}
	}
	dm.data.Name = "PIT"
	return &dm.data
}

func sortPlayers(players []game.PlayerID) []game.PlayerID {
	sort.Slice(players, func(i, j int) bool {
		return strings.Compare(string(players[i]), string(players[j])) < 0
	})
	return players
}

func (mg *GameStats) GetBattingData() *Data {
	var dm dataMaker
	for team, stats := range mg.TeamStats {
		if mg.filterExclude(team, stats) {
			continue
		}
		var players []game.PlayerID
		for player := range stats.Batting {
			players = append(players, player)
		}
		for _, player := range sortPlayers(players) {
			batting := stats.Batting[player]
			var m map[string]interface{}
			if err := mapstructure.Decode(batting, &m); err != nil {
				panic(err)
			}
			mg.adjustRowValues(len(batting.Games), team, batting.Player, m)
			dm.addRow(m)
		}
	}
	dm.data.Name = "BAT"
	return &dm.data
}

func (mg *GameStats) adjustRowValues(gameCount int, team string, player *game.Player, m map[string]interface{}) {
	delete(m, "Games")
	m["Games"] = gameCount
	if mg.OnlyTeam != "" {
		m["Name"] = player.NameOrNumber()
	} else {
		m["Name"] = fmt.Sprintf("%s/%s", team, player.NameOrNumber())
	}
	delete(m, "PlayerID")
	delete(m, "Number")
}

func (mg *GameStats) filterExclude(team string, stats *Stats) bool {
	if mg.OnlyTeam != "" &&
		!strings.HasPrefix(strings.ToLower(team), strings.ToLower(mg.OnlyTeam)) {
		return true
	}
	return false
}

func sortKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}