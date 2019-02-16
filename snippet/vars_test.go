package snippet

import "testing"

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
	for _, tt := range []struct {
		in   string
		out  string
		vars stringList
	}{
		{"lala ${0} ${1} ${2:bar} bla", "lala foo fuzz bar bla", []string{"foo", "fuzz"}},
		{"$0 - $1", "foo - bar", []string{"foo", "bar"}},
		{"foo ${a} $b ${}", "foo ${a} $b ${}", []string{"foo", "bar"}},
		{"{} [] () ]({", "{} [] () ]({", []string{"foo", "bar"}},
	} {
		t.Run(tt.in, func(t *testing.T) {
			v := expandVars(tt.in, tt.vars)
			if v != tt.out {
				t.Errorf("got %q, want %q", v, tt.out)
			}
		})
	}
}
