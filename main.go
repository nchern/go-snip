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

var (
	homeDir = os.Getenv("HOME")
)

func inHome(filename string) string {
	return path.Join(homeDir, filename)
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

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "%s recursively scans %q and reads all *.snip files found\n", os.Args[0], snippetsSrcRoot)
	}
}

type commandList []string

func (l commandList) String() string {
	return strings.Join(l, ", ")
}

var (
	commands = commandList{}

	snippetsSrcRoot = inHome(".vim")
	goSnipFile      = inHome(".go-snip")

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

	groups, err := snippet.LoadFromDir(snippetsSrcRoot)
	dieIf(err)

	err = save(groups)
	dieIf(err)

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
		err = snippets.PrintNames(os.Stdout)
		dieIf(err)
	} else if *cmd == cmdGroups {
		err = groups.PrintNames(os.Stdout)
		dieIf(err)
	} else {
		flag.Usage()
	}
}
