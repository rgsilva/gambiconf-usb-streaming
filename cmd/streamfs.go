// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program is the analogon of libfuse's hello.c, a a program that
// exposes a single file "file.txt" in the root directory.
package main

import (
	"fmt"
	"gambiconf/stream"
	"github.com/hanwen/go-fuse/v2/fs"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("%s <port> <raw disk image> <mount path>\n", os.Args[0])
		os.Exit(1)
	}

	var root fs.InodeEmbedder
	var mnt string

	bq := stream.NewBufferQueue()
	port, err := strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		panic(err)
	}

	receiver := stream.NewReceiver(bq, uint16(port))
	err = receiver.Start()
	if err != nil {
		panic(err)
	}

	root = &stream.Filesystem{
		DiskImage:   os.Args[2],
		FileSize:    int64((200 * 1024 * 1024) * 0.90),
		BufferQueue: bq,
	}
	mnt = os.Args[3]

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
			receiver.Stop()
		}
	}()

	server.Wait()
}
