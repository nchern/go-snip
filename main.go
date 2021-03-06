package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/nchern/go-snip/snippet"
)

const (
	customSnippetsRootVarName = "GOSNIP_SNIPPETS_ROOT"
)

var (
	homeDir                = os.Getenv("HOME")
	defaultSnippetsSrcRoot = inHome(".vim")
	goSnipFile             = inHome(".go-snip")
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "%s recursively scans %q and reads all *.snip files found\n", os.Args[0], snippetsSrcRoot())
	}
}

type commandList []string

func (l commandList) String() string {
	return strings.Join(l, ", ")
}

var (
	commands = commandList{}

	cmdLs     = c("ls")
	cmdShow   = c("show")
	cmdGroups = c("groups")

	cmd   = flag.String("cmd", "", fmt.Sprintf("One of: %s", commands))
	group = flag.String("g", "go", "Sets snippets group(snippet filename w/o .snip extension)")
)

func c(s string) string {
	commands = append(commands, s)
	return s
}

// go-snip -g=go -cmd=show func bar bazz
// go-snip -g=go -cmd=ls
// go-snip -g=go -cmd=groups
func main() {
	flag.Parse()

	groups, err := snippet.LoadFromDir(snippetsSrcRoot())
	dieIf(err)

	must(save(groups))

	snippets, found := groups[*group]
	if !found {
		dieIf(fmt.Errorf("group %q not found", *group))
	}

	if *cmd == cmdShow {
		name := flag.Arg(0)
		if s := snippets.Find(name); s != nil {
			fmt.Println(s.Render(flag.Args()[1:]))
			return
		}
		fmt.Println()
	} else if *cmd == cmdLs {
		must(snippets.PrintNames(os.Stdout))
	} else if *cmd == cmdGroups {
		must(groups.PrintNames(os.Stdout))
	} else {
		flag.Usage()
	}
}

func snippetsSrcRoot() string {
	if customRoot := strings.TrimSpace(os.Getenv(customSnippetsRootVarName)); customRoot != "" {
		return customRoot
	}
	return defaultSnippetsSrcRoot
}

func inHome(filename string) string {
	return path.Join(homeDir, filename)
}

func must(err error) {
	dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		log.Fatalf("// FATAL: %s\n", err)
	}
}

func save(m snippet.Groups) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(goSnipFile, body, 0644)
}
