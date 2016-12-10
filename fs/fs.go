package fs

/*
file structure:

mtpt/
	github/
		1000111/
			subject
			message
*/

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
)

type Article interface {
	ArticleID() uint32
	Subject() string
	Message() string
	Creation() time.Time
	LastMod() time.Time
}

type Service interface {
	Name() string
	List() ([]Article, error)
}

type Root struct {
	Dir
	services map[string]Service
}

func NewRoot() *Root {
	return &Root{
		Dir: NewDir(),
		srv: make(map[string]Service),
	}
}

func (root *Root) CreateService(srv Service) {
	root.services[srv.Name()] = srv
}

func (root *Root) ReadDir() []string {
	a := make([]string, 0, len(root.services))
	for name := range root.services {
		a = append(a, name)
	}
	return a
}

func (root *Root) MountAndServe(mtpt string) error {
	if err := mountAndServe(&dir); err != nil {
		return err
	}
}
