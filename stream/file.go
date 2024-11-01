package stream

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"io"
	"os"
	"syscall"
)

type File struct {
	fs.Inode

	fullPath       string
	original       *os.File
	stat           os.FileInfo
	offsetMin      int64
	offsetMax      int64
	readFromBuffer bool
	bufferQueue    *BufferQueue
}

// File has attributes
var _ = (fs.NodeGetattrer)((*File)(nil))

// File can be open
var _ = (fs.NodeOpener)((*File)(nil))

// File can be read
var _ = (fs.FileReader)((*File)(nil))

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size = uint64(f.stat.Size())
	return 0
}

func (f *File) Open(ctx context.Context, openFlags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	// disallow writes
	if openFlags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}

	// Return FOPEN_DIRECT_IO so content is not cached.
	return f, fuse.FOPEN_DIRECT_IO, 0
}

func (f *File) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	if off >= f.offsetMin && off <= f.offsetMax {
		if !f.readFromBuffer {
			f.readFromBuffer = f.bufferQueue.Has(128 * 1024)
		}

		if !f.readFromBuffer {
			// No data, return zeros.
			zeroes := make([]byte, len(dest))
			copy(dest, zeroes)
		} else {
			// Wait until we have the data.
			wait := f.bufferQueue.WaitFor(len(dest))
			<-wait

			// Read and return it.
			x := f.bufferQueue.Get(len(dest))
			copy(dest, x)
		}
		return fuse.ReadResultData(dest), 0
	}

	// Follow the default procedure.
	_, err := f.original.ReadAt(dest, off)
	if err != nil && err != io.EOF {
		return nil, syscall.EIO
	}

	return fuse.ReadResultData(dest), 0
}
