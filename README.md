[![Go Report Card](https://goreportcard.com/badge/github.com/nchern/go-snip)](https://goreportcard.com/report/github.com/nchern/go-snip)

# go-snip

Lightweight [Neosnippet](https://github.com/Shougo/neosnippet-snippets/tree/master/neosnippets) parser and processor. Just an experiment how to get snippets in vim w/o dealing with many various plugins.

The current implementation reads all the `*.snip` files in `~/.vim` dir recursively each time the command is invoked.

## Install 
```bash
go get github.com/nchern/go-snip/...
```

## Usage

### Command line

```bash
$ go-snip -g=go -cmd=show fori 'foo()' index MAX 
for index := 0; index < MAX; index++ {
foo()
}
```

### Plug it into vim

Let's say we want to have snippets for golang in vim. Install the command and then just add the following lines into your `.vimrc`:

```vim
"" Map GoSnip command to call the util
command! -range -nargs=* -complete=custom,ListSnippets GoSnip :<line1>,<line2>!go-snip -g=go -cmd=show <args>
"" Enable autocomplete for GoSnip command
:fun ListSnippets(A,L,P)
:    return system("go-snip -g=go -cmd=ls")
:endfun
```

Then in the command mode you can run a command like this: `:GoSnip fori 'foo()' index MAX`. It will insert the `for i ..` snippet at the cursor position.
.
## TODO
- [ ] read snippets from cache file at least for name listing
- [ ] configurable snippets search root
- [ ] multi search roots
