package snippet

import (
	"regexp"
	"strconv"
	"strings"
)

func parseVar(v string) (index int, defaultVal string, err error) {
	v = strings.TrimPrefix(v, "$")
	v = strings.Trim(v, "{}")
	toks := stringList(strings.Split(v, ":"))
	i, err := strconv.ParseInt(toks[0], 0, 64)
	if err != nil {
		return
	}
	index = int(i)
	defaultVal = toks.Get(1)
	return
}

func expandVar(v string, substitutions stringList) string {
	i, defaultVal, err := parseVar(v)
	if err != nil {
		return v
	}
	if val := substitutions.Get(i); val != "" {
		return val
	}
	return defaultVal
}

func expandVars(text string, substitutions stringList) string {
	re := regexp.MustCompile(`(\$\{.*?\}|\$\d+?)`)
	for _, v := range re.FindAllString(text, -1) {
		text = strings.Replace(text, v, expandVar(v, substitutions), -1)
	}
	return text
}
