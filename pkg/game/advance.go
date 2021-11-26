package game

import (
	"fmt"
	"regexp"
	"strings"
)

type Advance struct {
	Code               string
	From, To           string
	Out                bool  `yaml:",omitempty"`
	Fielders           []int `yaml:",omitempty,flow"`
	RunnerInterference bool  `yaml:",omitempty"`
	Implied            bool
	Runner             PlayerID
	WildPitch          bool
	PassedBall         bool
	*FieldingError     `yaml:",omitempty"`
}

type Advances map[string]*Advance

var advanceRegexp = regexp.MustCompile(`^([B123])([X-])([123H])(?:\(([^)]+)\))?$`)
var BaseNumber = map[string]int{
	"1": 0,
	"2": 1,
	"3": 2,
	"H": 3,
}
var PreviousBase = map[string]string{
	"H": "3",
	"3": "2",
	"2": "1",
}
var runnerNumber = map[string]int{
	"1": 0,
	"2": 1,
	"3": 2,
}

func (a *Advance) GoString() string {
	return a.Code
}

func parseAdvance(s string) (*Advance, error) {
	m := advanceRegexp.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("illegal advance code %s", s)
	}
	a := &Advance{
		Code: s,
		From: m[1],
		To:   m[3],
		Out:  m[2] == "X",
	}
	switch {
	case a.Out:
		if m[4] == "RINT" {
			a.RunnerInterference = true
		} else {
			for _, f := range m[4] {
				if f >= '1' && f <= '9' {
					a.Fielders = append(a.Fielders, int(f-'1')+1)
				} else {
					return nil, fmt.Errorf("illegal fielder %c for put out in advance code %s", f, s)
				}
			}
			if len(a.Fielders) == 0 {
				return nil, fmt.Errorf("no fielders for put out in advancde code %s", s)
			}
		}
	case m[4] == "WP":
		a.WildPitch = true
	case m[4] == "PB":
		a.PassedBall = true
	case m[4] != "":
		var err error
		a.FieldingError, err = parseFieldingError(m[4])
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func parseAdvances(advancesCode string, batter PlayerID, runners []PlayerID) (advances Advances, err error) {
	advances = make(Advances)
	if len(advancesCode) > 0 {
		for _, as := range strings.Split(advancesCode, ";") {
			var advance *Advance
			advance, err = parseAdvance(as)
			if err != nil {
				return
			}
			if advances[advance.From] != nil {
				err = fmt.Errorf("cannot advance %s twice in %s", advance.From, advancesCode)
				return
			}
			if advance.From == "B" {
				advance.Runner = batter
			} else {
				if runners == nil {
					err = fmt.Errorf("no runner to advance from %s at the start of a half-inning", advance.From)
					return
				}
				advance.Runner = runners[runnerNumber[advance.From]]
				if advance.Runner == "" {
					err = fmt.Errorf("no runner to advance from %s in %s", advance.From, advancesCode)
					return
				}
			}
			advances[advance.From] = advance
		}
	}
	return
}
