package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/morrowc/picam/client/client"
)

var (
	sAddr = flag.String("server_address", "127.0.0.1:443", "Remote image collection server address:port")
	id    = flag.String("id", "", "Identifier used for this image sender.")
	store = flag.String("store", "", "Directory where camera images are stored.")
)

func main() {
	flag.Parse()

	if *id == "" || *store == "" {
		fmt.Println("Provide an identifier and storage directory.")
		return
	}

	c, err := client.New(*sAddr, *id, *store)
	if err != nil {
		log.Fatalf("failed to create new picam client: %v", err)
	}

	// Start the watcher, and then monitor the channel.
	go c.Watcher()
	go func() {
		for fn := range c.Files {
			fmt.Printf("Watcher found a new file: %v\n", fn)
		}
	}()
}
