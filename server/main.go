package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang/glog"
	pgpb "github.com/morrowc/picam/proto/picam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	maxMsgSize = 5 * 1024 * 1024 // 5MB
)

var (
	port  = flag.Int("port", 9987, "Port upon which to listen.")
	store = flag.String("store", "/tmp/camstore", "Storage location for image files.")
)

// server holds connection/client information necessary to operate a gRPC server.
type server struct {
	store string
	port  int
	pgpb.UnimplementedPiCamServer
}

func new(store string, port int) *server {
	return &server{
		store: store,
		port:  port,
	}
}

func (s *server) SendImage(ctx context.Context, req *pgpb.Request) (*pgpb.Response, error) {
	id := req.GetIdentifier()
	img := req.GetImage()
	glog.Infof("Got request from ID: %s img size: %d", id, len(img))
	return &pgpb.Response{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		glog.Fatalf("failed to listen(): %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxMsgSize(maxMsgSize),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	server := new(*store, *port)

	pgpb.RegisterPiCamServer(s, server)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen&serve: %v", err)
	}

}
