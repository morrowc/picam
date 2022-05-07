package server

import (
	"fmt"

	"github.com/vladimirvivien/go4vl/v4l2/device"
)

type Server struct {
	dev *device.Device
	id  string
}

func New(dev, id string) (Server, error) {
	d, err := device.Open(dev)
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %v", err)
	}
	return &Server{
		dev: d,
		id:  id,
	}
}
