package main

import (
	"fmt"
	"gambiconf/timestamp"
	"github.com/hanwen/go-fuse/v2/fs"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%s <mount point>\n", os.Args[0])
		os.Exit(1)
	}

	var root fs.InodeEmbedder
	var mnt string

	root = &timestamp.Filesystem{}
	mnt = os.Args[1]

	opts := &fs.Options{}
	opts.Debug = false
	opts.AllowOther = true
	opts.ExplicitDataCacheControl = true

	server, err := fs.Mount(mnt, root, opts)
	if err != nil {
		panic(err)
	}

	// Listen for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			server.Unmount()
		}
	}()

	server.Wait()

}
