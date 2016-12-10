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

type Dir interface {
	Node
	Stat() *FileInfo
	ReadDir() []Dir
}

type Root struct {
	Node
	FileInfo
	services map[string]Service
}

func NewRoot() *Root {
	now := time.Now()
	return &Root{
		Node: NewNode(),
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

func (root *Root) Stat() *FileInfo {
	return &root.FileInfo
}

func (root *Root) ReadDir() []Dir {
	a := make([]Dir, 0, len(root.services))
	for name := range root.services {
		dir := &ServiceDir{
			Node: NewNode(),
			FileInfo: FileInfo{
				Name: name,
				Mode: os.ModeDir | 0755,
			},
		}
		a = append(a, dir)
	}
	return a
}

type ServiceDir struct {
	Node
	FileInfo
}

func (dir *ServiceDir) Stat() *FileInfo {
	return &dir.FileInfo
}

func (dir *ServiceDir) ReadDir() []Dir {
	return []Dir{}
}
