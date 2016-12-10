// +build linux

package fs

import (
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type Node interface {
	nodefs.Node
}

func NewNode() Node {
	return nodefs.NewDefaultNode()
}

func (f *FileInfo) FillAttr(out *fuse.Attr) {
	out.Mode = uint32(f.Mode & os.ModePerm)
	if f.IsDir() {
		out.Mode |= fuse.S_IFDIR
	}
	out.Size = uint64(f.Size)
	out.Atime = uint64(f.LastMod.Unix())
	out.Mtime = uint64(f.LastMod.Unix())
}

func (f *FileInfo) FillDirEntry(out *fuse.DirEntry) {
	out.Name = f.Name
	out.Mode = uint32(f.Mode & os.ModePerm)
	if f.IsDir() {
		out.Mode |= fuse.S_IFDIR
	}
}

func (root *Root) MountAndServe(mtpt string) error {
	opts := nodefs.Options{
		AttrTimeout:  time.Second,
		EntryTimeout: time.Second,
		Debug:        false,
	}
	s, _, err := nodefs.MountRoot(mtpt, root, &opts)
	if err != nil {
		return err
	}
	s.Serve()
	return nil
}

func (root *Root) Lookup(out *fuse.Attr, name string, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	return lookupName(root, name, out, ctx)
}

func (root *Root) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	root.FileInfo.FillAttr(out)
	return fuse.OK
}

func (root *Root) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return readDir(root)
}

func (dir *ServiceDir) Lookup(out *fuse.Attr, name string, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	return lookupName(dir, name, out, ctx)
}

func (dir *ServiceDir) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	dir.FileInfo.FillAttr(out)
	return fuse.OK
}

func (dir *ServiceDir) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return readDir(dir)
}

func (dir *TaskDir) Lookup(out *fuse.Attr, name string, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	return lookupName(dir, name, out, ctx)
}

func (dir *TaskDir) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	dir.FileInfo.FillAttr(out)
	return fuse.OK
}

func (dir *TaskDir) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return readDir(dir)
}

func (t *Text) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	t.FileInfo.FillAttr(out)
	return fuse.OK
}

func (t *Text) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return nil, fuse.EINVAL
}

func (t *Text) Open(flags uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	p, err := t.ReadFile()
	if err != nil {
		return nil, fuse.EIO
	}
	return nodefs.NewDataFile(p), fuse.OK
}

func lookupName(dir Dir, name string, out *fuse.Attr, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	_, status := readDir(dir)
	if status != fuse.OK {
		return nil, status
	}
	c := dir.Inode().GetChild(name)
	if c == nil {
		return nil, fuse.ENOENT
	}
	status = c.Node().GetAttr(out, nil, ctx)
	if status != fuse.OK {
		return nil, status
	}
	return c, fuse.OK
}

func readDir(dir Dir) ([]fuse.DirEntry, fuse.Status) {
	p := dir.Inode()
	kids, err := dir.ReadDir()
	if err != nil {
		return nil, fuse.EIO
	}
	a := make([]fuse.DirEntry, len(kids))
	for i, kid := range kids {
		info := kid.Stat()
		if p.GetChild(info.Name) == nil {
			p.NewChild(info.Name, info.IsDir(), kid)
		}
		info.FillDirEntry(&a[i])
	}
	return a, fuse.OK
}
