// Package main provides the tool to walks through the file system directories and print relative to
// $HOME path of every matched directory pattern.

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime/debug"
	"sync"
)

func main() {
	debug.SetGCPercent(-1)
	if err := doMain(); err != nil {
		log.Panicln(err)
	}
}

type opts struct {
	homeDir     string
	rootDirs    []string
	patternDirs []string
	staticDirs  []string
}

func parseOpts() opts {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	n := len(os.Args)
	if n < 3 {
		panic(errors.New("few number of arguments"))
	}
	opts := opts{homeDir: usr.HomeDir}
	sep := 0
	for ; sep < n && os.Args[sep] != "--patterns"; sep++ {
	}

	opts.rootDirs = os.Args[1:sep]
	sep = min(sep+1, n)
	i := sep
	for ; sep < n && os.Args[sep] != "--static"; sep++ {
	}

	opts.patternDirs = os.Args[i:sep]
	sep = min(sep+1, n)
	opts.staticDirs = os.Args[sep:n]

	return opts
}

func doMain() error {
	opts := parseOpts()

	var wg sync.WaitGroup
	wg.Add(len(opts.rootDirs))

	for _, dirName := range opts.rootDirs {
		go func(dirName string) {
			defer wg.Done()
			_ = filepath.WalkDir(dirName, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !d.IsDir() {
					return nil
				}

				i := 0
				n := len(opts.patternDirs)
				for ; i < n; i++ {
					dir := opts.patternDirs[i]
					name := filepath.Join(path, dir)
					if _, err := os.Stat(name); err == nil {
						break
					}
				}
				if i == n {
					return nil
				}

				dir, err := filepath.Rel(opts.homeDir, path)
				if err == nil {
					fmt.Fprintln(os.Stdout, dir)
				}

				return filepath.SkipDir
			})
		}(dirName)
	}

	wg.Wait()

	for _, dir := range opts.staticDirs {
		dir, err := filepath.Rel(opts.homeDir, dir)
		if err == nil {
			fmt.Fprintln(os.Stdout, dir)
		}
	}

	return nil
}
