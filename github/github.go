package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v74/github"
	"github.com/lufia/taskfs/fs"
	"golang.org/x/oauth2"
)

type Comment struct {
	seq     int
	comment *github.IssueComment
}

func NewComment(seq int, comment *github.IssueComment) *Comment {
	return &Comment{
		seq:     seq,
		comment: comment,
	}
}

func (p *Comment) Key() string {
	return fmt.Sprintf("%d", p.seq)
}

func (p *Comment) Message() string {
	return *p.comment.Body
}

func (p *Comment) Creation() time.Time {
	return p.comment.CreatedAt.Time
}

func (p *Comment) LastMod() time.Time {
	return p.comment.UpdatedAt.Time
}

type Issue struct {
	issue *github.Issue
	svc   *Service
}

func (p *Issue) Key() string {
	owner := p.repositoryOwner()
	repo := p.repositoryName()
	return fmt.Sprintf("%s@%s#%d", repo, owner, p.Number())
}

func (p *Issue) Subject() string {
	return *p.issue.Title
}

func (p *Issue) Message() string {
	return *p.issue.Body
}

func (p *Issue) PermaLink() string {
	return *p.issue.HTMLURL
}

func (p *Issue) Creation() time.Time {
	return p.issue.CreatedAt.Time
}

func (p *Issue) LastMod() time.Time {
	return p.issue.UpdatedAt.Time
}

func (p *Issue) Comments() (a []fs.Comment, err error) {
	var buf []*github.IssueComment
	page := 0
	for {
		var b []*github.IssueComment
		b, page, err = p.fetchComments(page)
		if err != nil {
			return
		}
		buf = append(buf, b...)
		if page == 0 {
			break
		}
	}
	a = make([]fs.Comment, len(buf))
	for i, v := range buf {
		a[i] = NewComment(i+1, v)
	}
	return a, nil
}

func (p *Issue) repositoryOwner() string {
	owner := *p.issue.Repository.Owner.Login
	if org := p.issue.Repository.Organization; org != nil {
		owner = *org.Login
	}
	return owner
}

func (p *Issue) repositoryName() string {
	return *p.issue.Repository.Name
}

func (p *Issue) Number() int {
	return *p.issue.Number
}

func (p *Issue) fetchComments(page int) ([]*github.IssueComment, int, error) {
	ctx := context.Background()
	owner := p.repositoryOwner()
	repo := p.repositoryName()
	n := p.Number()
	var opt github.IssueListCommentsOptions
	opt.Page = page
	b, resp, err := p.svc.c.Issues.ListComments(ctx, owner, repo, n, &opt)
	if err != nil {
		return nil, 0, err
	}
	return b, resp.NextPage, nil
}

type Config struct {
	BaseURL string
	Token   string
}

func (c *Config) authorizedClient() *http.Client {
	if c.Token == "" {
		return nil
	}
	token := &oauth2.Token{
		AccessToken: c.Token,
	}
	s := oauth2.StaticTokenSource(token)
	return oauth2.NewClient(oauth2.NoContext, s)
}

type Service struct {
	c    *github.Client
	name string
}

func NewService(config *Config) (*Service, error) {
	var client *http.Client
	if config.Token != "" {
		client = config.authorizedClient()
	}
	c := github.NewClient(client)
	name := "github.com"
	if config.BaseURL != "" {
		u, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, err
		}
		c.BaseURL = u
		name = u.Host
	}
	return &Service{c: c, name: name}, nil
}

func (p *Service) Name() string {
	return p.name
}

func (p *Service) List() ([]fs.Task, error) {
	var a []fs.Task
	var opt github.IssueListOptions
	ctx := context.Background()
	for {
		b, resp, err := p.c.Issues.List(ctx, true, &opt)
		if err != nil {
			return nil, err
		}
		a = p.appendIssues(a, b)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return a, nil
}

func (p *Service) appendIssues(a []fs.Task, b []*github.Issue) []fs.Task {
	for _, v := range b {
		a = append(a, &Issue{issue: v, svc: p})
	}
	return a
}
