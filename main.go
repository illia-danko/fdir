package main

import (
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
}

func parseOpts() opts {
	n := len(os.Args)
	sep := 0
	for ; sep < n && os.Args[sep] != "--patterns"; sep++ {
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return opts{
		homeDir:     usr.HomeDir,
		rootDirs:    os.Args[1:sep],
		patternDirs: os.Args[sep+1:],
	}
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

	return nil
}
