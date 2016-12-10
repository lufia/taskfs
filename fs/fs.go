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
	Key() string
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
	now := time.Now()
	dirs := make([]Dir, 0, len(root.services))
	for name, svc := range root.services {
		dir := &ServiceDir{
			Node: NewNode(),
			FileInfo: FileInfo{
				Name:     name,
				Mode:     os.ModeDir | 0755,
				Creation: now,
				LastMod:  now,
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
	svc   Service
	cache []Dir
}

func (dir *ServiceDir) Stat() *FileInfo {
	return &dir.FileInfo
}

func (dir *ServiceDir) ReadDir() ([]Dir, error) {
	if dir.cache != nil {
		return dir.cache, nil
	}
	a, err := dir.svc.List()
	if err != nil {
		return nil, err
	}
	dirs := make([]Dir, len(a))
	for i, task := range a {
		dirs[i] = &TaskDir{
			Node: NewNode(),
			FileInfo: FileInfo{
				Name:     task.Key(),
				Mode:     os.ModeDir | 0755,
				Creation: task.Creation(),
				LastMod:  task.LastMod(),
			},
			task: task,
		}
	}
	dir.cache = dirs
	return dirs, nil
}

type TaskDir struct {
	Node
	FileInfo
	task Task
}

func (dir *TaskDir) Stat() *FileInfo {
	return &dir.FileInfo
}

func (dir *TaskDir) ReadDir() ([]Dir, error) {
	return []Dir{}, nil
}
