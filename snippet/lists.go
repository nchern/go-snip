package snippet

import (
	"fmt"
	"io"
)

type stringList []string

func (l stringList) Get(i int) string {
	if i < len(l) {
		return l[i]
	}
	return ""
}

type list []*snippet

func (l list) add(s *snippet) list {
	if s == nil || s.body == nil {
		return l
	}
	return append(l, s)
}

func (l list) Find(name string) *snippet {
	// TODO: handle alias
	for _, s := range l {
		if s.name == name {
			return s
		}
	}
	return nil
}

func (l list) PrintNames(w io.Writer) error {
	for _, s := range l {
		if _, err := fmt.Fprintln(w, s.name); err != nil {
			return err
		}
	}
	return nil
}
