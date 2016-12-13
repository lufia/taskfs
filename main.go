package main

import (
	"flag"
	"log"

	"github.com/lufia/taskfs/fs"
	"github.com/lufia/taskfs/github"
	"github.com/lufia/taskfs/gitlab"
)

var (
	accessToken   = flag.String("t", "", "access token")
	baseURL       = flag.String("url", "", "endpoint url")
	glAccessToken = flag.String("gitlab.t", "", "access token")
	glBaseURL     = flag.String("gitlab.url", "", "endpoint url")
	debug         = flag.Bool("d", false, "turn on debug print")
)

func main() {
	flag.Parse()
	s, err := github.NewService(&github.Config{
		BaseURL: *baseURL,
		Token:   *accessToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	g, err := gitlab.NewService(&gitlab.Config{
		BaseURL: *glBaseURL,
		Token:   *glAccessToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	root := fs.NewRoot()
	root.CreateService(s)
	root.CreateService(g)
	if err := root.MountAndServe("x", *debug); err != nil {
		log.Fatal(err)
	}
}
