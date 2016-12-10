// +build linux

package fs

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type Dir struct {
	nodefs.Node
}

func NewDir() Dir {
	return Dir{
		Node: nodefs.NewDefaultNode(),
	}
}

type File struct {
	nodefs.Node
}

func NewFile(p []byte) File {
	return File{
		Node: nodefs.NewDataFile(p),
	}
}

func (dir *ServiceDir) MountAndServe(mtpt string) error {
	opts := nodefs.Options{
		AttrTimeout:  time.Second,
		EntryTimeout: time.Second,
		Debug:        false,
	}
	s, _, err := nodefs.MountRoot(mtpt, dir, opts)
	if err != nil {
		return err
	}
	s.Serve()
}
