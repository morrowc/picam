// Package client builds a web-cam image capture client
// to speak to the remote picam server.
package client

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/gidoBOSSftw5731/log"
	"google.golang.org/grpc"

	pgpb "github.com/morrowc/picam/proto/picam"
)

const (
	maxMsgSize = 1024 * 1000 * 10 // 10mB
)

// Client holds all of the information about the running client.
type Client struct {
	client   pgpb.PiCamClient
	id       string
	srvAddr  string
	store    string
	files    chan string
	ImgCount int64
}

func New(srvAddr, id, store string) (*Client, error) {
	conn, err := grpc.Dial(
		srvAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxMsgSize)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make new connection: %v", err)
	}

	return &Client{
		client:   pgpb.NewPiCamClient(conn),
		id:       id,
		srvAddr:  srvAddr,
		store:    store,
		files:    make(chan string, 10),
		ImgCount: 0,
	}, nil
}

// Watcher starts a watch process on the store, sending write events to the channel.
func (c *Client) Watcher() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Error creating the file watcher: %v", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					c.files <- event.Name
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Infof("error: %v", err)
			}
		}
	}()
	return nil
}

// SendImage, Send an image to the remote server.
func (c *Client) SendImage(ctx context.Context, img []byte) error {
	// Build a request and send it to the server.
	req := &pgpb.Request{
		Identifier: c.id,
		Image:      img,
	}

	resp, err := c.client.SendImage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send image: %v - ", err, resp.GetError())
	}
	c.ImgCount++
	log.Infof("Successfully uploaded image, now %d sent.", c.ImgCount)
	return nil
}
