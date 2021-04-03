package snippet

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	snippetText = `
snippet foo
alias f
abbr foo ...
    foobar

# comment
snippet multiline
alias m
abbr multiline ...
    bar

	fuzzbuzz


`
)

func TestTestShouldStringOK(t *testing.T) {
	const given = `
snippet foo
alias f
abbr foo ...
    foobar

snippet multiline
alias m
abbr multiline ...
    bar
        fuzzbuzz
`
	underTest, err := parse(bytes.NewBufferString(given))
	assert.NoError(t, err)
	assert.Equal(t, "foobar\n", underTest[0].String())
	assert.Equal(t, "bar\n    fuzzbuzz", underTest[1].String())
}

func TestShouldParse(t *testing.T) {
	var tests = []struct {
		name     string
		given    string
		expected list
	}{
		{"empty", "", []*snippet{}},
		{
			"two_with_multiline",
			snippetText,
			[]*snippet{
				{
					name:  "foo",
					alias: "f",
					abbr:  "foo ...",
					body: []string{
						"    foobar",
						"",
					},
					bodyMinIndent: 4,
				},
				{
					name:  "multiline",
					alias: "m",
					abbr:  "multiline ...",
					body: []string{
						"    bar",
						"",
						"\tfuzzbuzz",
						"",
						"",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actualList, err := parse(bytes.NewBufferString(tt.given))

			assert.NoError(t, err)
			assert.Len(t, actualList, len(tt.expected))
			for i, actual := range actualList {
				assert.Equal(t, tt.expected[i], actual)
			}
		})
	}
}

func TestShouldParseFail(t *testing.T) {
	expected := errors.New("boom")
	r := &errReader{expected}
	res, err := parse(r)
	assert.Nil(t, res)
	assert.Equal(t, expected, err)
}

func TestShouldFindSnippet(t *testing.T) {
	underTest := list([]*snippet{
		{name: "foo"},
		{name: "bar", abbr: "first bar"},
		{name: "bar", alias: "b"},
		{name: "buzz", alias: "b"},
		{name: "fuzz"},
	})

	var tests = []struct {
		name          string
		expectedIndex int
		given         string
	}{
		{"find by name", 3, "buzz"},
		{"find by name - searches first match", 1, "bar"},
		{"find by alias - searches first match", 2, "b"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := underTest.Find(tt.given)
			assert.Equal(t, underTest[tt.expectedIndex], actual)
		})
	}

	assert.Nil(t, underTest.Find("not_existent"))
}

type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (n int, err error) { return 0, r.err }
