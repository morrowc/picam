// Package is a client side service implementation to collect images
// from a fsnotify watched directory, and deliver to a picam service server.
package main

import (
	"context"
	"flag"
	"log"
	"sync"

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
		log.Fatal("Failed to provide ID or SERVER flag values.")
	}

	// Create a client service, start fswatching the store directory.
	c, err := client.New(*server, *id, *store)
	if err != nil {
		log.Fatalf("failed to create client service: %v", err)
	}
	if err := c.Watcher(); err != nil {
		log.Fatalf("failed to start the fsnotify watcher: %v", err)
	}

	// Start a simple waitgroup to watch/wait on the image sending
	// loop to return, and allow the program to exit.
	wg := sync.WaitGroup()

	wg.Add(1)
	// Run a goroutine to just loop watching for images to send.
	go func() {
		defer wg.Done()
		ctx := context.Background()
		for {
			fn := <-c.files
			img, err := os.ReadFile(fn)
			if err != nil {
				log.Errorf("failed to read the stored image file(%s): %v", fn, err)
				return
			}
			if err := c.SendImage(ctx, img); err != nil {
				log.Errorf("failed to send image: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	log.Infof("Finished waiting for images, send: %d images.", c.ImgCount)

}
