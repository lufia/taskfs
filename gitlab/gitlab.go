package gitlab

import (
	"fmt"
	"net/url"
	"time"

	"github.com/lufia/taskfs/fs"
	"github.com/xanzy/go-gitlab"
)

type Comment struct {
	num  int
	note *gitlab.Note
}

func (p *Comment) Key() string {
	return fmt.Sprintf("%d", p.num)
}

func (p *Comment) Message() string {
	return p.note.Body
}

func (p *Comment) Creation() time.Time {
	return *p.note.CreatedAt
}

func (p *Comment) LastMod() time.Time {
	return *p.note.UpdatedAt
}

type Issue struct {
	issue *gitlab.Issue
	proj  *gitlab.Project
	svc   *Service
}

func (p *Issue) Key() string {
	owner := p.proj.Namespace.Name
	repo := p.proj.Name
	return fmt.Sprintf("%s@%s#%d", owner, repo, p.issue.IID)
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
	var buf []*gitlab.Note
	page := 0
	for {
		var b []*gitlab.Note
		b, page, err = p.fetchNotes(page)
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
		a[i] = &Comment{num: i + 1, note: v}
	}
	return a, nil
}

func (p *Issue) fetchNotes(page int) ([]*gitlab.Note, int, error) {
	pid := p.issue.ProjectID
	n := p.issue.ID
	var opt gitlab.ListIssueNotesOptions
	opt.Page = page
	b, resp, err := p.svc.c.Notes.ListIssueNotes(pid, n, &opt)
	if err != nil {
		return nil, 0, err
	}
	return b, resp.NextPage, nil
}

type Config struct {
	BaseURL string
	Token   string
}

type Service struct {
	c        *gitlab.Client
	name     string
	projects map[int]*gitlab.Project
}

func NewService(config *Config) (*Service, error) {
	u, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}
	c := gitlab.NewClient(nil, config.Token)
	c.SetBaseURL(config.BaseURL)
	svc := &Service{
		c:        c,
		name:     u.Host,
		projects: make(map[int]*gitlab.Project),
	}
	return svc, nil
}

func (p *Service) Name() string {
	return p.name
}

func (p *Service) List() ([]fs.Task, error) {
	var a []fs.Task
	var opt gitlab.ListIssuesOptions
	for {
		b, resp, err := p.c.Issues.ListIssues(&opt)
		if err != nil {
			return nil, err
		}
		a, err = p.convertAppendIssues(a, b)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return a, nil
}

func (p *Service) convertAppendIssues(a []fs.Task, b []*gitlab.Issue) ([]fs.Task, error) {
	for _, v := range b {
		task, err := p.fetchTask(v)
		if err != nil {
			return nil, err
		}
		a = append(a, task)
	}
	return a, nil
}

func (p *Service) fetchTask(v *gitlab.Issue) (task fs.Task, err error) {
	proj := p.projects[v.ProjectID]
	if proj == nil {
		proj, _, err = p.c.Projects.GetProject(v.ProjectID)
		if err != nil {
			return
		}
		p.projects[v.ProjectID] = proj
	}
	return &Issue{issue: v, proj: proj, svc: p}, nil
}
