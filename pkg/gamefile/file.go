package gamefile

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/participle/v2/lexer"
)

type Numbers string
type Position = lexer.Position

type File struct {
	Path          string
	Properties    map[string]string
	PropertyPos   map[string]Position
	VisitorEvents []*Event
	HomeEvents    []*Event

	PropertyList []*Property   `parser:"@@*"`
	TeamEvents   []*TeamEvents `parser:"@@*"`
}

type TeamEvents struct {
	Pos           Position
	HomeOrVisitor string   `parser:"@('visitorplays' | 'homeplays') (NL|EOF)"`
	Events        []*Event `parser:"@@*"`
}

type Property struct {
	Pos   Position
	Key   string `parser:"@Key"`
	Value string `parser:"@Value (NL+|EOF)"`
}

type Event struct {
	Pos         Position
	Alternative *Alternative    `parser:"'alt' @@ (NL|EOF)"`
	Pitcher     string          `parser:"| ('pitcher'|'pitching') @Token (NL|EOF)"`
	RAdjRunner  Numbers         `parser:"| 'radj' @Token"`
	RAdjBase    string          `parser:"      @Token (NL|EOF)"`
	Score       string          `parser:"| 'score' @Token (NL|EOF)"`
	Final       string          `parser:"| 'final' @Token (NL|EOF)"`
	Sub         *LineupChange   `parser:"| @@"`
	Play        *ActualPlay     `parser:"| @@"`
	Afters      []*LineupChange `parser:"   @@*"`
	Comment     string          `parser:"   @Comment? (NL|EOF)"`
	Empty       bool            `parser:"| @NL"`
}

type LineupChange struct {
	CourtesyRunner *string `parser:"'cr' @Token"`
	Conference     *bool   `parser:"| @'conf'"`
	SubEnter       string  `parser:"| 'sub' @Token"`
	SubExit        string  `parser:"|     'for' @Token"`
	HSubEnter      string  `parser:"| 'hsub' @Token"`
	HSubExit       string  `parser:"      'for' @Token (NL|EOF)"`
	VSubEnter      string  `parser:"| 'vsub' @Token"`
	VSubExit       string  `parser:"      'for' @Token (NL|EOF)"`
}

type Play interface {
	GetPos() Position
	GetCode() string
	GetAdvances() []string
}

type ActualPlay struct {
	Pos                      Position
	ContinuedPlateAppearance bool     `parser:"((@'...')"`
	PlateAppearance          Numbers  `parser:" | (@PA"`
	Batter                   string   `parser:"    @Token))"`
	PitchSequence            string   `parser:" @Token"`
	Code                     string   `parser:" @Token"`
	Advances                 []string `parser:" @Advance*"`
}

var _ Play = (*ActualPlay)(nil)

type Alternative struct {
	Pos      Position
	Code     string   `parser:"@Token"`
	Advances []string `parser:"@Advance*"`
	Comment  string   `parser:"  @Comment?"`
}

var _ Play = (*Alternative)(nil)

func (n *Numbers) UnmarshalText(dat []byte) error {
	*n = Numbers(strings.TrimRight(string(dat), " \t"))
	return nil
}

func (n Numbers) String() string {
	return string(n)
}

func (n Numbers) Int() int {
	i, _ := strconv.Atoi(n.String())
	return i
}

func (f *File) Parse(r io.Reader) error {
	return nil
}

func (p *ActualPlay) normalize() {
	if p == nil {
		return
	}
	for i, adv := range p.Advances {
		p.Advances[i] = strings.ToUpper(adv)
	}
	p.Code = strings.ToUpper(p.Code)
	p.PitchSequence = strings.ToUpper(p.PitchSequence)
}

func (a *Alternative) normalize() {
	if a == nil {
		return
	}
	for i, adv := range a.Advances {
		a.Advances[i] = strings.ToUpper(adv)
	}
	a.Code = strings.ToUpper(a.Code)
}

func (f *File) Validate() error {
	f.Properties = make(map[string]string)
	f.PropertyPos = make(map[string]Position)
	for _, prop := range f.PropertyList {
		f.Properties[prop.Key] = prop.Value
		f.PropertyPos[prop.Key] = prop.Pos
	}
	for _, te := range f.TeamEvents {
		for _, event := range te.Events {
			// make codes upper code
			event.Play.normalize()
			event.Alternative.normalize()
		}
		switch te.HomeOrVisitor {
		case "homeplays":
			if f.HomeEvents != nil {
				return fmt.Errorf("%s: duplicate homeplays section", te.Pos)
			}
			f.HomeEvents = append(f.HomeEvents, te.Events...)
		case "visitorplays":
			if f.VisitorEvents != nil {
				return fmt.Errorf("%s: duplicate visitorplays section", te.Pos)
			}
			f.VisitorEvents = append(f.VisitorEvents, te.Events...)
		}
	}
	f.setPlateAppearances(f.HomeEvents)
	f.setPlateAppearances(f.VisitorEvents)
	return nil
}

func (f *File) setPlateAppearances(events []*Event) {
	var pa Numbers
	for _, event := range events {
		if event.Play != nil {
			if event.Play.ContinuedPlateAppearance {
				event.Play.PlateAppearance = pa
			} else {
				pa = event.Play.PlateAppearance
			}
		}
	}
}

func (f *File) Write(w io.Writer) {
	printed := map[string]bool{}
	for _, name := range []string{"date", "game", "visitor", "visitorid", "home", "homeid", "start", "timelimit", "tournament", "league"} {
		val := f.Properties[name]
		printed[name] = true
		if val != "" {
			fmt.Fprintf(w, "%s: %s\n", name, val)
		}
	}
	var names []string
	for name := range f.Properties {
		if !printed[name] {
			printed[name] = true
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		val := f.Properties[name]
		fmt.Fprintf(w, "%s: %s\n", name, val)
	}
	fmt.Fprintln(w, "---")
	f.writeEvents(w, "visitorplays", f.VisitorEvents)
	f.writeEvents(w, "homeplays", f.HomeEvents)
}

func (f *File) writeEvents(w io.Writer, name string, events []*Event) {
	if events == nil {
		return
	}
	fmt.Fprintf(w, "%s\n", name)
	var pa int
	for _, event := range events {
		switch {
		case event.Play != nil:
			play := event.Play
			if !play.ContinuedPlateAppearance {
				if i := play.PlateAppearance.Int(); i != 0 {
					pa = i
				} else {
					pa += 1
				}
				fmt.Fprintf(w, "%d %s ", pa, play.Batter)
			} else {
				fmt.Fprintf(w, "  ... ")
			}
			fmt.Fprintf(w, "%s ", play.PitchSequence)
			f.writeCodeAdvancesComment(w, play.Code, play.Advances, event.Afters, event.Comment)
		case event.Alternative != nil:
			alt := event.Alternative
			fmt.Fprintf(w, "  alt ")
			f.writeCodeAdvancesComment(w, alt.Code, alt.Advances, nil, alt.Comment)
		case event.Pitcher != "":
			fmt.Fprintf(w, "pitching %s\n", event.Pitcher)
		case event.RAdjBase != "":
			fmt.Fprintf(w, "radj %s %s\n", event.RAdjRunner, event.RAdjBase)
		case event.Score != "":
			fmt.Fprintf(w, "score %s\n", event.Score)
		case event.Final != "":
			fmt.Fprintf(w, "final %s\n", event.Final)
		}
	}
}

func (f *File) writeCodeAdvancesComment(w io.Writer, code string, advances []string, afters []*LineupChange, comment string) {
	fmt.Fprintf(w, "%s", code)
	for _, adv := range advances {
		fmt.Fprintf(w, " %s", adv)
	}
	for _, aft := range afters {
		if aft.Conference != nil {
			fmt.Fprint(w, " conf")
		}
		if aft.CourtesyRunner != nil {
			fmt.Fprintf(w, " cr %s", *aft.CourtesyRunner)
		}
	}
	if comment != "" {
		fmt.Fprintf(w, " : %s", comment)
	}
	fmt.Fprintln(w)
}

const GameDateFormat = "1/2/2006"

func (f *File) GetGameDate() (time.Time, error) {
	d := f.Properties["date"]
	t, err := time.Parse("1/2/06", d)
	if err != nil {
		t, err = time.Parse(GameDateFormat, d)
	}
	if err != nil {
		return t, fmt.Errorf("%s: can't parse date: %w", f.PropertyPos["date"], err)
	}
	return t, nil
}

func (p *ActualPlay) GetPos() Position {
	return p.Pos
}

func (p *ActualPlay) GetCode() string {
	return p.Code
}

func (p *ActualPlay) GetAdvances() []string {
	return p.Advances
}

func (a *Alternative) GetPos() Position {
	return a.Pos
}

func (a *Alternative) GetCode() string {
	return a.Code
}

func (a *Alternative) GetAdvances() []string {
	return a.Advances
}

func (a *Alternative) GetComment() string {
	return a.Comment
}
