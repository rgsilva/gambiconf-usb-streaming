package corrupt

import (
	"context"
	"fmt"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"os"
	"path"
	"syscall"
)

type (
	Filesystem struct {
		fs.Inode

		OriginalDir string
	}
)

// Filesystem has its files added on OnAdd
var _ = (fs.NodeOnAdder)((*Filesystem)(nil))

// Filesystem can have its files opened
var _ = (fs.NodeOpener)((*Filesystem)(nil))

func (r *Filesystem) OnAdd(ctx context.Context) {
	existingFiles, err := os.ReadDir(r.OriginalDir)
	if err != nil {
		panic(err)
	}

	for _, existingFile := range existingFiles {
		filename := existingFile.Name()
		fmt.Printf("Adding: %s\n", filename)

		stat, err := existingFile.Info()
		if err != nil {
			panic(err)
		}

		fullPath := path.Join(r.OriginalDir, filename)
		ch := r.NewPersistentInode(ctx, &File{
			fullPath: fullPath,
			stat:     stat,
		}, fs.StableAttr{})
		r.AddChild(filename, ch, true)
	}
}

func (r *Filesystem) Open(ctx context.Context, openFlags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	// disallow writes
	if fuseFlags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}

	// Return FOPEN_DIRECT_IO so content is not cached.
	return fh, fuse.FOPEN_DIRECT_IO, 0
}
