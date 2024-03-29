// Package is a client side service implementation to collect images
// from a fsnotify watched directory, and deliver to a picam service server.
package main

import (
	"context"
	"flag"
	"os"
	"path"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/morrowc/picam/client/client"
)

var (
	store  = flag.String("store", "/tmp/camstore", "Storage directory for images.")
	id     = flag.String("id", "", "Identifier of this instance of the service.")
	server = flag.String("server", "", "Server address: hostname:port")
)

func main() {
	flag.Parse()
	if *id == "" || *server == "" {
		glog.Fatal("Failed to provide ID or SERVER flag values.")
	}

	// Create a client service, start fswatching the store directory.
	c, err := client.New(*server, *id, *store)
	if err != nil {
		glog.Fatalf("failed to create client service: %v", err)
	}
	go func() error {
		if err := c.Watcher(); err != nil {
			glog.Fatalf("failed to start the fsnotify watcher: %v", err)
			return err
		}
		return nil
	}()

	// Start a simple waitgroup to watch/wait on the image sending
	// loop to return, and allow the program to exit.
	var wg sync.WaitGroup

	wg.Add(1)
	// Run a goroutine to just loop watching for images to send.
	go func() {
		defer wg.Done()
		ctx := context.Background()
		for {
			fn := <-c.Files
			if fn == "" {
				time.Sleep(2 * time.Second)
				continue
			}
			// Delay a second to let file writes finish.
			time.Sleep(1 * time.Second)
			img, err := os.ReadFile(fn)
			if err != nil {
				glog.Errorf("failed to read the stored image file(%s): %v", fn, err)
				return
			}
			glog.Infof("Extracted file: %v which is %d bytes in size.", fn, len(img))
			if err := c.SendImage(ctx, path.Base(fn), img); err != nil {
				glog.Errorf("failed to send image: %v", err)
			}
		}
	}()
	wg.Wait()
	glog.Infof("Finished waiting for images, send: %d images.", c.ImgCount)
}
