package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

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
	return []fs.Task{}, nil
}
