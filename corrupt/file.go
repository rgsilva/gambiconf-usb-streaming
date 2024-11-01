package corrupt

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"io"
	"math/rand"
	"os"
	"syscall"
)

type File struct {
	fs.Inode

	fullPath string
	original *os.File
	stat     os.FileInfo
}

// File has attributes
var _ = (fs.NodeGetattrer)((*File)(nil))

// File can be read
var _ = (fs.FileReader)((*File)(nil))

// File can be open
var _ = (fs.NodeOpener)((*File)(nil))

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size = uint64(f.stat.Size())
	return 0
}

func (f *File) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n, err := f.original.ReadAt(dest, off)
	if err != nil && err != io.EOF {
		return nil, syscall.EIO
	}

	if n == 0 && err == io.EOF {
		// If nothing was read and EOF was reached, return empty data
		return fuse.ReadResultData(nil), 0
	}

	for i := 0; i < n; i++ {
		if rand.Intn(5) == 0 {
			// Pick a random printable ASCII character and replace the byte
			dest[i] = byte(rand.Intn(126-32+1) + 32)
		}
	}

	return fuse.ReadResultData(dest[:n]), 0
}

func (f *File) Open(ctx context.Context, openFlags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	// disallow writes
	if openFlags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}

	fileDescriptor, err := os.Open(f.fullPath)
	if err != nil {
		panic(err)
	}

	f.original = fileDescriptor

	// Return FOPEN_DIRECT_IO so content is not cached.
	return f, fuse.FOPEN_DIRECT_IO, 0
}
