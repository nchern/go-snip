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

type Groups map[string]List

func (g Groups) PrintNames(w io.Writer) error {
	for k := range g {
		if _, err := fmt.Fprintln(w, k); err != nil {
			return err
		}
	}
	return nil
}

type Snippet struct {
	Name string

	Abbr string

	Alias string

	Body []string
}

func (s *Snippet) String() string {
	return strings.Join(s.Body, "\n")
}

func (s *Snippet) Render(vals []string) string {
	return expandVars(s.String(), vals)
}

type stringList []string

func (l stringList) Get(i int) string {
	if i < len(l) {
		return l[i]
	}
	return ""
}

type List []*Snippet

func (l List) add(s *Snippet) List {
	if s == nil || s.Body == nil {
		return l
	}
	return append(l, s)
}

func (l List) Find(name string) *Snippet {
	// TODO: handle alias
	for _, s := range l {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (l List) PrintNames(w io.Writer) error {
	for _, s := range l {
		if _, err := fmt.Fprintln(w, s.Name); err != nil {
			return err
		}
	}
	return nil
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

func parse(reader io.Reader) (List, error) {
	res := List{}
	var current *Snippet

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
			current = &Snippet{Name: tokens[1]}
		} else if current != nil && snippetLine(line).IsAbbr() {
			current.Abbr = line
		} else if current != nil && snippetLine(line).IsAlias() {
			tokens := strings.Fields(line)
			current.Alias = stringList(tokens).Get(1)
		} else if current != nil {
			current.Body = append(current.Body, line)
		}
	}
	res = res.add(current)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func parseFile(filename string) (List, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parse(file)
}

func LoadFromDir(rootDir string) (Groups, error) {
	// this implementation OVERWRITES group for files with the same names

	res := Groups{}
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ext) {
			return nil
		}
		_, filename := filepath.Split(info.Name())
		key := strings.TrimSuffix(filename, ext)
		var err error
		res[key], err = parseFile(path)
		return err
	}); err != nil {
		return nil, err
	}
	return res, nil
}
