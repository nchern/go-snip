# go-snip

Lightweight [Neosnippet](https://github.com/Shougo/neosnippet-snippets/tree/master/neosnippets) parser and processor. Just an experiment how to get snippets in vim w/o dealing with many various plugins.

### Install 
```bash
go get github.com/nchern/go-snip/...
```

### Usage
```bash
$ go-snip -g=go -cmd=show fori 'foo()' index MAX 
for index := 0; index < MAX; index++ {
foo()
}
```

## How to plug it into vim

Let's say we want to have snippets for golang. Install the command and then just add the following lines into your `.vimrc`

```vim
"" Map GoSnip command to call the util
command! -range -nargs=* -complete=custom,ListSnippets GoSnip :<line1>,<line2>!go-snip -g=go -cmd=show <args>
"" Enable autocomplete for GoSnip command
:fun ListSnippets(A,L,P)
:    return system("go-snip -g=go -cmd=ls")
:endfun
```

## TODO
- [ ] read snippets from cache file at least for name listing
- [ ] configurable snippets search root
- [ ] multi search roots
