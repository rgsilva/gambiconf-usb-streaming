package timestamp

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
	"time"
)

type File struct {
	fs.Inode

	content []byte
}

// File can be open
var _ = (fs.NodeOpener)((*File)(nil))

// File can be read
var _ = (fs.FileReader)((*File)(nil))

func (f *File) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	end := off + int64(len(dest))
	if end > int64(len(f.content)) {
		end = int64(len(f.content))
	}

	// We could copy to the `dest` buffer, but since we have a
	// []byte already, return that.
	return fuse.ReadResultData(f.content[off:end]), 0
}

func (f *File) Open(ctx context.Context, openFlags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	// disallow writes
	if openFlags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}

	// capture open time
	now := time.Now().Format(time.StampNano) + "\n"
	f.content = []byte(now)

	// Return FOPEN_DIRECT_IO so content is not cached.
	return f, fuse.FOPEN_DIRECT_IO, 0
}
