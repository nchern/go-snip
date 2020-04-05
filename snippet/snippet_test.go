package snippet

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const snippetText = `
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

func TestShouldParse(t *testing.T) {
	var tests = []struct {
		name     string
		given    string
		expected list
	}{
		{"empty", "", []*snippet{}},
		{"two_with_multiline", snippetText, []*snippet{
			&snippet{name: "foo", alias: "f", abbr: "foo ...", body: []string{"foobar"}},
			&snippet{name: "multiline", alias: "m", abbr: "multiline ...", body: []string{"bar", "fuzzbuzz"}}}},
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

type errReader struct {
	err error
}

func (r *errReader) Read(p []byte) (n int, err error) { return 0, r.err }
