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

type Task interface {
	TaskID() uint32
	Subject() string
	Message() string
	Creation() time.Time
	LastMod() time.Time
}

type Service interface {
	Name() string
	List() ([]Task, error)
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
	ReadDir() ([]Dir, error)
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

func (root *Root) ReadDir() ([]Dir, error) {
	dirs := make([]Dir, 0, len(root.services))
	for name, svc := range root.services {
		dir := &ServiceDir{
			Node: NewNode(),
			FileInfo: FileInfo{
				Name: name,
				Mode: os.ModeDir | 0755,
			},
			svc: svc,
		}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}

type ServiceDir struct {
	Node
	FileInfo
	svc Service
}

func (dir *ServiceDir) Stat() *FileInfo {
	return &dir.FileInfo
}

func (dir *ServiceDir) ReadDir() ([]Dir, error) {
	a, err := dir.svc.List()
	if err != nil {
		return nil, err
	}
	dirs := make([]Dir, len(a))
	for i, _ := range a {
		dirs[i] = &TaskDir{}
	}
	return []Dir{}, nil
}

type TaskDir struct {
	Node
	FileInfo
}

func (dir *TaskDir) Stat() *FileInfo {
	return &dir.FileInfo
}

func (dir *TaskDir) ReadDir() ([]Dir, error) {
	return []Dir{}, nil
}
