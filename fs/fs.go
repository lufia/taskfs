package fs

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
)

type Server interface {
}

type Article interface {
	ArticleID() uint32
	Subject() string
	Message() string
	Creation() time.Time
	LastMod() time.Time
}

type Dir interface {
}

type File interface {
	Stat(attr *fuse.Attr)
	Data() []byte
}
