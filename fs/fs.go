package fs

/*
file structure:

mtpt/
	github/
		1000111/
			subject
			message
			1
			2
			3...
*/

import (
	"errors"
	"os"
	"strings"
	"time"
)

type Comment interface {
	Key() string
	Message() string
	Creation() time.Time
	LastMod() time.Time
}

type Task interface {
	Key() string
	Subject() string
	Message() string
	PermaLink() string
	Creation() time.Time
	LastMod() time.Time
	Comments() ([]Comment, error)
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
	registers map[string]func(token, url string) (Service, error)
	services  map[string]Service
}

func NewRoot() *Root {
	now := time.Now()
	return &Root{
		Node: NewNode(),
		FileInfo: FileInfo{
			Mode:     os.ModeDir | 0755,
			Creation: now,
			LastMod:  now,
		},
		registers: make(map[string]func(token, url string) (Service, error)),
		services:  make(map[string]Service),
	}
}

func (root *Root) RegisterService(kind string, fn func(token, url string) (Service, error)) {
	if _, ok := root.registers[kind]; ok {
		panic("duplicate service register: " + kind)
	}
	root.registers[kind] = fn
}

func (root *Root) Stat() *FileInfo {
	return &root.FileInfo
}

func (root *Root) ReadDir() ([]Dir, error) {
	now := time.Now()
	dirs := make([]Dir, 0, len(root.services)+1) // +1: ctl file
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
	dirs = append(dirs, &Ctl{
		Node: NewNode(),
		FileInfo: FileInfo{
			Name:     "ctl",
			Mode:     0644,
			Creation: now,
			LastMod:  now,
		},
		Commands: map[string]func(args ...string) error{
			"add": root.addService,
		},
	})
	return dirs, nil
}

func (root *Root) addService(args ...string) error {
	var kind, token, url string
	switch len(args) {
	case 3:
		url = args[2]
		fallthrough
	case 2:
		token = args[1]
		kind = args[0]
		register := root.registers[kind]
		if register == nil {
			return errors.New("unsupported service type: " + kind)
		}
		srv, err := register(token, url)
		if err != nil {
			return err
		}
		root.services[srv.Name()] = srv
		return nil
	default:
		return errors.New("invalid add command")
	}
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
	dirs := make([]Dir, len(a)+1)
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
	now := time.Now()
	dirs[len(dirs)-1] = &Ctl{
		Node: NewNode(),
		FileInfo: FileInfo{
			Name:     "ctl",
			Mode:     0644,
			Creation: now,
			LastMod:  now,
		},
		Commands: map[string]func(args ...string) error{
			"refresh": dir.refreshCache,
		},
	}
	dir.cache = dirs
	return dirs, nil
}

func (dir *ServiceDir) refreshCache(args ...string) error {
	dir.cache = nil
	return nil
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
	a, err := dir.task.Comments()
	if err != nil {
		return nil, err
	}
	kids := make([]Dir, 0, len(a)+3)
	kids = append(kids, dir.newText("subject", dir.task.Subject()))
	kids = append(kids, dir.newText("message", dir.task.Message()))
	kids = append(kids, dir.newText("url", dir.task.PermaLink()))
	for _, c := range a {
		kids = append(kids, NewCommentText(c))
	}
	dir.files = kids
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

type CommentText struct {
	Node
	FileInfo
	data []byte
}

func NewCommentText(c Comment) *CommentText {
	data := []byte(c.Message())
	return &CommentText{
		Node: NewNode(),
		FileInfo: FileInfo{
			Name:     c.Key(),
			Size:     int64(len(data)),
			Mode:     0644,
			Creation: c.Creation(),
			LastMod:  c.LastMod(),
		},
		data: data,
	}
}

func (t *CommentText) Stat() *FileInfo {
	return &t.FileInfo
}

func (t *CommentText) ReadDir() ([]Dir, error) {
	return nil, errProtocol
}

func (t *CommentText) ReadFile() ([]byte, error) {
	return t.data, nil
}

type Ctl struct {
	Node
	FileInfo
	Commands map[string]func(args ...string) error
}

func (ctl *Ctl) Stat() *FileInfo {
	return &ctl.FileInfo
}

func (ctl *Ctl) ReadDir() ([]Dir, error) {
	return nil, errProtocol
}

func (ctl *Ctl) ReadFile() ([]byte, error) {
	return []byte{}, nil
}

func (ctl *Ctl) WriteFile(p []byte) error {
	s := string(p)
	cmds := strings.Split(s, "\n")
	for _, cmd := range cmds {
		a := strings.Fields(cmd)
		if len(a) == 0 {
			return nil
		}
		fn, ok := ctl.Commands[a[0]]
		if !ok {
			return errors.New("unknown control command")
		}
		if err := fn(a[1:]...); err != nil {
			return err
		}
	}
	return nil
}
