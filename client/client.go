// Package client builds a web-cam image capture client
// to speak to the remote picam server.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/rjeczalik/notify"
	"google.golang.org/grpc"

	pgpb "github.com/morrowc/picam/proto/picam"
)

const (
	maxMsgSize = 1024 * 1000 * 10 // 10mB
)

// Client holds all of the information about the running client.
type Client struct {
	client   pgpb.PiCamClient
	Id       string
	srvAddr  string
	store    string
	Files    chan string
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
		Id:       id,
		srvAddr:  srvAddr,
		store:    store,
		Files:    make(chan string, 1),
		ImgCount: 0,
	}, nil
}

// Watcher starts a watch process on the store, sending write events to the channel.
func (c *Client) Watcher() error {
	e := make(chan notify.EventInfo, 1)
	glog.Infof("Watcher starting for: %s", c.store)
	if err := notify.Watch(c.store, e, notify.InCreate, notify.InMovedTo); err != nil {
		return fmt.Errorf("Error creating the file watcher: %v", err)
	}
	// Close the watcher and the channel of files when this function returns.
	defer notify.Stop(e)
	defer close(e)

	// Start a watching goroutine.
	for {
		glog.Info("Looping on events")
		switch ei := <-e; ei.Event() {
		case notify.InCreate:
			c.Files <- ei.Path()
			glog.Infof("Create event for path: %s", ei.Path())
		case notify.InMovedTo:
			c.Files <- ei.Path()
			glog.Infof("InMovedTo event for path: %s", ei.Path())
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// SendImage, Send an image to the remote server.
func (c *Client) SendImage(ctx context.Context, fn string, img []byte) error {
	// Build a request and send it to the server.
	req := &pgpb.Request{
		Identifier: c.Id,
		Image:      img,
		Filename:   fn,
	}

	resp, err := c.client.SendImage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send image: %v - %v", err, resp.GetError())
	}
	c.ImgCount++
	glog.Infof("Successfully uploaded image, now %d sent.", c.ImgCount)
	return nil
}
