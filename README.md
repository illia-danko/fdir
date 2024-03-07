The tool walks through the file system directories and print relative to $HOME
path of every matched directory pattern.

# Usage

```bash
go run main.go "$HOME/github.com" "$HOME/gitlab.com" --patterns .git .hg .bzr .svn
```

Output:

```
github.com/clearloop/leetcode-cli
github.com/elixir-lsp/elixir-ls
...
```

In the example above, `fdir` tool seeks recursively into `$HOME/github.com` and
`$HOME/github.com` directories until it finds any of `.git`, `.hd`, `bzr` or
`.svn` folder. When the find has occurred it prints the parent folder of the
matched pattern.
