package snippet

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	ext = ".snip"
)

type group struct {
	list

	sourceFilename string
}

// LoadFromDir traverses given root dir, parses snippets and returns populated snippet group collection
func LoadFromDir(rootDir string) (Groups, error) {
	// this implementation OVERWRITES group for files with the same names

	res := Groups{}
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ext) {
			return nil
		}
		_, filename := filepath.Split(info.Name())
		key := strings.TrimSuffix(filename, ext)
		snippets, err := parseFile(path)
		if err != nil {
			return err
		}
		res[key] = &group{
			list:           snippets,
			sourceFilename: filepath.Join(path, info.Name()),
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

// Groups represents a mapping if group names to snippet group data structure
type Groups map[string]*group

// PrintNames outputs group names along with their the corresponding file names to a given writer
func (g Groups) PrintNames(w io.Writer) error {
	for k, v := range g {
		if _, err := fmt.Fprintf(w, "%s\t(%s)\n", k, v.sourceFilename); err != nil {
			return err
		}
	}
	return nil
}

type snippet struct {
	name string

	abbr string

	alias string

	body []string

	bodyMinIndent int
}

func (s *snippet) String() string {
	lines := []string{}
	for _, l := range s.body {
		if s.bodyMinIndent > 0 && strings.TrimSpace(l) != "" {
			l = l[s.bodyMinIndent:]
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func (s *snippet) Render(vals []string) string {
	return expandVars(s.String(), vals)
}

type snippetLine string

func (l snippetLine) Fields() []string {
	return strings.Fields(string(l))
}

func (l snippetLine) IsAbbr() bool {
	return strings.HasPrefix(string(l), "abbr ")
}

func (l snippetLine) IsSnippet() bool {
	return strings.HasPrefix(string(l), "snippet ")
}

func (l snippetLine) IsAlias() bool {
	return strings.HasPrefix(string(l), "alias ")
}

func (l snippetLine) IsComment() bool {
	return strings.HasPrefix(string(l), "#")
}

func (l snippetLine) IsBlank() bool {
	return strings.TrimSpace(string(l)) == ""
}

func parse(reader io.Reader) (list, error) {
	res := list{}
	var current *snippet

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := snippetLine(strings.TrimSpace(scanner.Text()))

		if line.IsComment() {
			continue
		} else if line.IsSnippet() {
			res = res.add(current)
			tokens := line.Fields()
			if len(tokens) < 2 {
				current = nil
				continue
			}
			current = &snippet{name: tokens[1], bodyMinIndent: int(math.MaxInt64)}
		} else if current != nil && line.IsAbbr() {
			current.abbr = strings.TrimSpace(strings.TrimPrefix(string(line), line.Fields()[0]))
		} else if current != nil && line.IsAlias() {
			tokens := line.Fields()
			current.alias = stringList(tokens).Get(1)
		} else if current != nil {
			ind := len(scanner.Text()) - len(strings.TrimLeft(scanner.Text(), " "))
			if ind < current.bodyMinIndent && line != "" {
				current.bodyMinIndent = ind
			}
			current.body = append(current.body, scanner.Text())
		}
	}
	res = res.add(current)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func parseFile(filename string) (list, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parse(file)
}
