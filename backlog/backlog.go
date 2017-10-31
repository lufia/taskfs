package backlog

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	backlog "github.com/griffin-stewie/go-backlog"
	"github.com/lufia/taskfs/fs"
)

type Issue struct {
	issue *backlog.Issue
	svc   *Service
}

func (p *Issue) Key() string {
	return *p.issue.IssueKey
}

func (p *Issue) Subject() string {
	return *p.issue.Summary
}

func (p *Issue) Message() string {
	return *p.issue.Description
}

func (p *Issue) PermaLink() string {
	return fmt.Sprintf("https://%s/view/%s", p.svc.name, *p.issue.IssueKey)
}

func (p *Issue) Creation() time.Time {
	return *p.issue.Created
}

func (p *Issue) LastMod() time.Time {
	return *p.issue.Updated
}

func (p *Issue) Comments() ([]fs.Comment, error) {
	return []fs.Comment{}, nil
}

type Config struct {
	BaseURL string
	APIKey  string
}

type Service struct {
	c      *backlog.Client
	name   string
	userID int
}

var (
	errMissingURL = errors.New("base url is missing")
)

func NewService(config *Config) (*Service, error) {
	if config.BaseURL == "" {
		return nil, errMissingURL
	}
	u, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}
	name := u.Host

	c := backlog.NewClient(u, config.APIKey)
	user, err := c.Myself()
	if err != nil {
		return nil, err
	}
	return &Service{c: c, name: name, userID: *user.ID}, nil
}

func (p *Service) Name() string {
	return p.name
}

func (p *Service) List() ([]fs.Task, error) {
	l, err := p.c.IssuesWithOption(&backlog.IssuesOption{
		AssigneeIDs: []int{p.userID},
		Statuses: []backlog.IssueStatus{
			backlog.Open,
			backlog.InProgress,
			backlog.Resolved,
		},
	})
	if err != nil {
		return nil, err
	}
	a := make([]fs.Task, len(l))
	for i, v := range l {
		a[i] = &Issue{issue: v, svc: p}
	}
	return a, nil
}
