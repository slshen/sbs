package gamefile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type YAMLParser struct {
	err error
}

func ParseYAMLFile(path string) (*File, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := &YAMLParser{}
	return p.parse(path, dat)
}

func (p *YAMLParser) parse(path string, dat []byte) (*File, error) {
	var m map[string]interface{}
	if err := yaml.Unmarshal(dat, &m); err != nil {
		return nil, err
	}
	f := &File{Path: path}
	pos := lexer.Position{Filename: path, Line: 1}
	for key, value := range m {
		switch key {
		case "homeplays":
			f.HomeEvents = p.parseYAMLEvents(p.guessPosition(path, "homeplays:", dat), value)
		case "visitorplays":
			f.VisitorEvents = p.parseYAMLEvents(p.guessPosition(path, "visitorplays:", dat), value)
		default:
			if val := p.toString(value); val != "" {
				f.PropertyList = append(f.PropertyList, &Property{
					Pos:   pos,
					Key:   key,
					Value: val,
				})
			}
			continue
		}
	}
	if err := f.Validate(); err != nil {
		p.err = multierror.Append(p.err, err)
	}
	return f, p.err
}

func (p *YAMLParser) guessPosition(path, key string, dat []byte) lexer.Position {
	s := bufio.NewScanner(bytes.NewReader(dat))
	i := 1
	for s.Scan() {
		if s.Text() == key {
			break
		}
		i++
	}
	return lexer.Position{
		Filename: path,
		Line:     i,
	}
}

func (p *YAMLParser) parseYAMLEvents(pos lexer.Position, value interface{}) (events []*Event) {
	plays, ok := value.([]interface{})
	if !ok {
		return
	}
	pa := 1
	for _, s := range plays {
		pos.Line++
		code := p.toString(s)
		if code == "" {
			return
		}
		parts := strings.Split(code, ",")
		switch parts[0] {
		case "pitcher":
			events = append(events, &Event{
				Pos:     pos,
				Pitcher: p.getPart(parts, 1),
			})
		case "inn":
			if len(parts) > 2 {
				events = append(events, &Event{
					Pos:   pos,
					Score: p.getPart(parts, 2),
				})
			}
		case "final":
			events = append(events, &Event{
				Pos:   pos,
				Final: p.getPart(parts, 1),
			})
		case "radj":
			events = append(events, &Event{
				Pos:        pos,
				RAdjRunner: Numbers(p.getPart(parts, 1)),
				RAdjBase:   p.getPart(parts, 2),
			})
		case "err":
			// ignore
		default:
			play := &ActualPlay{
				Pos:             pos,
				Batter:          p.parseBatter(p.getPart(parts, 0)),
				PitchSequence:   p.getPart(parts, 1),
				PlateAppearance: Numbers(fmt.Sprintf("%d", pa)),
			}
			code := p.getPart(parts, 2)
			dot := strings.IndexRune(code, '.')
			if dot > 0 {
				play.Code = code[0:dot]
				if dot+1 < len(code) {
					play.Advances = strings.Split(code[dot+1:], ";")
				}
			} else {
				play.Code = code
			}
			events = append(events, &Event{
				Pos:     pos,
				Play:    play,
				Comment: p.getPart(parts, 3),
			})
			if p.isPAComplete(code) {
				pa++
			}
		}
	}
	return
}

func (p *YAMLParser) isPAComplete(code string) bool {
	if strings.HasPrefix(code, "SB") ||
		strings.HasPrefix(code, "WP") ||
		strings.HasPrefix(code, "PB") ||
		strings.HasPrefix(code, "CS") ||
		strings.HasPrefix(code, "PO") ||
		code == "NP" ||
		strings.HasPrefix(code, "FLE") {
		return false
	}
	return true
}

func (p *YAMLParser) parseBatter(s string) string {
	// the yaml format allows letters at the start of a batter
	// but the gamefile format only allows digits, so remove the
	// letters
	m := regexp.MustCompile(`[a-z]*([0-9]+)`).FindStringSubmatch(s)
	if m != nil {
		return m[1]
	}
	return "000"
}

func (p *YAMLParser) getPart(parts []string, i int) string {
	if i < len(parts) {
		return parts[i]
	}
	return ""
}

func (p *YAMLParser) toString(s interface{}) string {
	if s == nil {
		return ""
	}
	switch v := s.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case []interface{}:
		if len(v) == 1 {
			return fmt.Sprintf("%s", v[0])
		}
		return ""
	default:
		return ""
	}
}
