package main

import (
	"flag"
	"log"

	"github.com/lufia/taskfs/github"
)

var (
	accessToken = flag.String("t", "", "access token")
	baseURL     = flag.String("url", "", "endpoint url")
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
	a, err := s.List()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range a {
		log.Println(v.ArticleID(), v.Subject())
		log.Println(v.Message())
	}
}
