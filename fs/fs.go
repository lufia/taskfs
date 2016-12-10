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
	"os"
	"time"
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

type FileInfo struct {
	Name     string
	Size     int64
	Mode     os.FileMode
	Creation time.Time
	LastMod  time.Time
}

func (f *FileInfo) IsDir() bool {
	return f.Mode&os.ModeDir != 0
}

type Root struct {
	Dir
	FileInfo
	services map[string]Service
}

func NewRoot() *Root {
	now := time.Now()
	return &Root{
		Dir: NewDir(),
		FileInfo: FileInfo{
			Mode:     os.ModeDir | 0777,
			Creation: now,
			LastMod:  now,
		},
		services: make(map[string]Service),
	}
}

func (root *Root) CreateService(srv Service) {
	root.services[srv.Name()] = srv
}

func (root *Root) ReadDir() []*FileInfo {
	a := make([]*FileInfo, 0, len(root.services))
	for name := range root.services {
		a = append(a, &FileInfo{
			Name: name,
			Size: 0,
			Mode: os.ModeDir | 0777,
		})
	}
	return a
}
