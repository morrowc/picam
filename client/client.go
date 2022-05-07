// Package client builds a web-cam image capture client
// to speak to the remote picam server.
package client

import (
	"context"
	"fmt"

	"github.com/vladimirvivien/go4vl/v4l2"
	"github.com/vladimirvivien/go4vl/v4l2/device"
	"google.golang.org/grpc"
)

const (
	maxMsgSize = 1024 * 1000 * 10 // 10mB
	format     = v4l2.PixelFmtMJPEG
)

// Client holds all of the information about the running client.
type Client struct {
	dev      *device.Device
	conn     *grpc.ClientConn
	id       string
	h        int
	w        int
	srvAddr  string
	imgCount int64
}

func New(dev, srvAddr, id, string, h, w int) (Server, error) {
	d, err := device.Open(dev)
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %v", err)
	}

	return &Server{
		dev:      d,
		conn:     nil,
		id:       id,
		h:        h,
		w:        w,
		srvAddr:  srvAddr,
		imgCount: 0,
	}
}

func (c *Client) newConn(ctx context.Context) error {
	conn, err := grpc.Dial(
		c.srvAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxMsgSize)),
	)
	if err != nil {
		return fmt.Errorf("failed to make new connection: %v", err)
	}
	c.conn = conn
	return nil
}

// func (c *Client) Setup
