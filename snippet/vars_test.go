package snippet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVar(t *testing.T) {
	for _, tt := range []struct {
		in              string
		expectedIndex   int
		expectedDefault string
	}{
		{"$13", 13, ""},
		{"${0}", 0, ""},
		{"${2}", 2, ""},
		{"${1:fn}", 1, "fn"},
		{"${10:someUrl}", 10, "someUrl"},
	} {
		t.Run(tt.in, func(t *testing.T) {
			i, dv, err := parseVar(tt.in)
			if i != tt.expectedIndex {
				t.Errorf("got %d, want %d", i, tt.expectedIndex)
			}
			if dv != tt.expectedDefault {
				t.Errorf("got %q, want %q", dv, tt.expectedDefault)
			}
			if err != nil {
				t.Errorf("got %s", err)
			}
		})
	}
}

func TestExpandVar(t *testing.T) {
	for _, tt := range []struct {
		in  string
		out string
	}{
		{"$0", "foo"},
		{"${0}", "foo"},
		{"${1:bar}", "bar"},
	} {
		t.Run(tt.in, func(t *testing.T) {
			v := expandVar(tt.in, []string{"foo"})
			if v != tt.out {
				t.Errorf("got %q, want %q", v, tt.out)
			}
		})
	}
}

func TestExpandVars(t *testing.T) {
	const multilineWithComment = `
	type ${0:Interface} interface {
         ${1:/* TODO: add methods */}
     }`
	const multilineWithCommentExpanded = `
	type Foo interface {
         /* TODO: add methods */
     }`

	for _, tt := range []struct {
		name     string
		given    string
		expected string
		vars     stringList
	}{
		{"two provided one default", "lala ${0} ${1} ${2:bar} bla", "lala foo fuzz bar bla", []string{"foo", "fuzz"}},
		{"expanded both", "$0 - $1", "foo - bar", []string{"foo", "bar"}},
		{"vars in curly brackets", "foo ${a} $b ${}", "foo ${a} $b ${}", []string{"foo", "bar"}},
		{"different brackets does not break", "{} [] () ]({", "{} [] () ]({", []string{"foo", "bar"}},
		{"multiline with comment", multilineWithComment, multilineWithCommentExpanded, []string{"Foo"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := expandVars(tt.given, tt.vars)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
