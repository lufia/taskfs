package main

import (
	"flag"
	"log"

	"github.com/lufia/taskfs/backlog"
	"github.com/lufia/taskfs/fs"
	"github.com/lufia/taskfs/github"
	"github.com/lufia/taskfs/gitlab"
)

var (
	debug = flag.Bool("d", false, "turn on debug print")

	mtpt = "/mnt/taskfs"
)

func main() {
	flag.Parse()
	root := fs.NewRoot()
	root.RegisterService("github", func(token, url string) (fs.Service, error) {
		return github.NewService(&github.Config{
			BaseURL: url,
			Token:   token,
		})
	})
	root.RegisterService("gitlab", func(token, url string) (fs.Service, error) {
		return gitlab.NewService(&gitlab.Config{
			BaseURL: url,
			Token:   token,
		})
	})
	root.RegisterService("backlog", func(token, url string) (fs.Service, error) {
		return backlog.NewService(&backlog.Config{
			BaseURL: url,
			APIKey:  token,
		})
	})
	if flag.NArg() > 0 {
		mtpt = flag.Arg(0)
	}
	if err := root.MountAndServe(mtpt, *debug); err != nil {
		log.Fatal(err)
	}
}
