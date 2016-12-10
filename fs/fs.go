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
	"errors"
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
	ReadFile() ([]byte, error)
}

var errProtocol = errors.New("protocol botch")

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

func (*Root) ReadFile() ([]byte, error) {
	return nil, errProtocol
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

func (*ServiceDir) ReadFile() ([]byte, error) {
	return nil, errProtocol
}

type TaskDir struct {
	Node
	FileInfo
	task  Task
	files []Dir
}

func (dir *TaskDir) Stat() *FileInfo {
	return &dir.FileInfo
}

func (dir *TaskDir) ReadDir() ([]Dir, error) {
	if dir.files != nil {
		return dir.files, nil
	}
	dir.files = []Dir{
		dir.newText("subject", dir.task.Subject()),
		dir.newText("message", dir.task.Message()),
	}
	return dir.files, nil
}

func (*TaskDir) ReadFile() ([]byte, error) {
	return nil, errProtocol
}

func (dir *TaskDir) newText(name, s string) *Text {
	data := []byte(s)
	return &Text{
		Node: NewNode(),
		FileInfo: FileInfo{
			Name:     name,
			Size:     int64(len(data)),
			Mode:     0644,
			Creation: dir.task.Creation(),
			LastMod:  dir.task.LastMod(),
		},
		data: data,
	}
}

type Text struct {
	Node
	FileInfo
	data []byte
}

func (t *Text) Stat() *FileInfo {
	return &t.FileInfo
}

func (t *Text) ReadDir() ([]Dir, error) {
	// this method isn't going to be called.
	return nil, errProtocol
}

func (t *Text) ReadFile() ([]byte, error) {
	return t.data, nil
}
