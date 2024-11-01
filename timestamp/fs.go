package timestamp

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
	"time"
)

type (
	Filesystem struct {
		fs.Inode
	}
)

// Filesystem can have files be opened.
var _ = (fs.NodeOpener)((*Filesystem)(nil))

// Filesystem has its files added on OnAdd
var _ = (fs.NodeOnAdder)((*Filesystem)(nil))

func (r *Filesystem) Open(ctx context.Context, openFlags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	// disallow writes
	if fuseFlags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}

	// capture open time
	now := time.Now().Format(time.StampNano) + "\n"
	fh = &File{
		content: []byte(now),
	}

	// Return FOPEN_DIRECT_IO so content is not cached.
	return fh, fuse.FOPEN_DIRECT_IO, 0
}

func (r *Filesystem) OnAdd(ctx context.Context) {
	ch := r.NewPersistentInode(
		ctx,
		&File{},
		fs.StableAttr{})
	r.AddChild("clock", ch, true)
}
