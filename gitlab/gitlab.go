package gitlab

import (
	"fmt"
	"time"

	"github.com/lufia/taskfs/fs"
	"github.com/xanzy/go-gitlab"
)

type Issue struct {
	issue *gitlab.Issue
	svc   *Service
}

func (p *Issue) Key() string {
	// TODO: implement
	//group
	//project
	return fmt.Sprintf("%d", p.issue.ID)
}

func (p *Issue) Subject() string {
	return p.issue.Title
}

func (p *Issue) Message() string {
	return p.issue.Description
}

func (p *Issue) PermaLink() string {
	return p.issue.WebURL
}

func (p *Issue) Creation() time.Time {
	return *p.issue.CreatedAt
}

func (p *Issue) LastMod() time.Time {
	return *p.issue.UpdatedAt
}

func (p *Issue) Comments() (a []fs.Comment, err error) {
	// TODO: implement
	return
}

type Config struct {
	BaseURL string
	Token   string
}

type Service struct {
	c *gitlab.Client
}

func NewService(config *Config) (*Service, error) {
	c := gitlab.NewClient(nil, config.Token)
	c.SetBaseURL(config.BaseURL)
	return &Service{c: c}, nil
}

func (p *Service) Name() string {
	return "gitlab"
}

func (p *Service) List() ([]fs.Task, error) {
	var a []fs.Task
	var opt gitlab.ListIssuesOptions
	for {
		b, resp, err := p.c.Issues.ListIssues(&opt)
		if err != nil {
			return nil, err
		}
		a = p.appendIssues(a, b)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return a, nil
}

func (p *Service) appendIssues(a []fs.Task, b []*gitlab.Issue) []fs.Task {
	for _, v := range b {
		a = append(a, &Issue{issue: v, svc: p})
	}
	return a
}
