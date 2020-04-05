package snippet

import (
	"bufio"
	"fmt"
	"io"
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

// PrintNames outputs group names along with thier the corresponding file names to a given writer
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
}

func (s *snippet) String() string {
	return strings.Join(s.body, "\n")
}

func (s *snippet) Render(vals []string) string {
	return expandVars(s.String(), vals)
}

type snippetLine string

func (l snippetLine) IsAbbr() bool {
	return strings.HasPrefix(string(l), "abbr")
}

func (l snippetLine) IsSnippet() bool {
	return strings.HasPrefix(string(l), "snippet")
}

func (l snippetLine) IsAlias() bool {
	return strings.HasPrefix(string(l), "alias")
}

func (l snippetLine) IsCommentOrBlank() bool {
	return strings.HasPrefix(string(l), "#") || l == ""
}

func parse(reader io.Reader) (list, error) {
	res := list{}
	var current *snippet

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if snippetLine(line).IsCommentOrBlank() {
			continue
		} else if snippetLine(line).IsSnippet() {
			res = res.add(current)
			tokens := strings.Fields(line)
			if len(tokens) < 2 {
				current = nil
				continue
			}
			current = &snippet{name: tokens[1]}
		} else if current != nil && snippetLine(line).IsAbbr() {
			current.abbr = line
		} else if current != nil && snippetLine(line).IsAlias() {
			tokens := strings.Fields(line)
			current.alias = stringList(tokens).Get(1)
		} else if current != nil {
			current.body = append(current.body, line)
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
