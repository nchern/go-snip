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

type Groups map[string]*group

func (g Groups) PrintNames(w io.Writer) error {
	for k, v := range g {
		if _, err := fmt.Fprintf(w, "%s\t(%s)\n", k, v.sourceFilename); err != nil {
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

type list []*Snippet

func (l list) add(s *Snippet) list {
	if s == nil || s.Body == nil {
		return l
	}
	return append(l, s)
}

func (l list) Find(name string) *Snippet {
	// TODO: handle alias
	for _, s := range l {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (l list) PrintNames(w io.Writer) error {
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

func parse(reader io.Reader) (list, error) {
	res := list{}
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

func parseFile(filename string) (list, error) {
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
