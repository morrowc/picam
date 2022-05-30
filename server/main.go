package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang/glog"
	pb "github.com/morrowc/picam/proto/picam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/prototext"
)

const (
	maxMsgSize = 5 * 1024 * 1024 // 5MB
)

var (
	port   = flag.Int("port", 9987, "Port upon which to listen.")
	store  = flag.String("store", "/tmp/camstore", "Storage location for image files.")
	config = flag.String("config", "", "Configuration for the server.")
)

// server holds connection/client information necessary to operate a gRPC server.
type server struct {
	store  string
	port   int
	config *pb.Config
	pb.UnimplementedPiCamServer
}

func new(store string, port int, config *pb.Config) *server {
	return &server{
		store:  store,
		port:   port,
		config: config,
	}
}

func (s *server) SendImage(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	id := req.GetIdentifier()
	img := req.GetImage()
	glog.Infof("Got request from ID: %s img size: %d", id, len(img))
	return &pb.Response{}, nil
}

// readConfig reads the server configuration proto from disk.
func readConfig(fn string) (*pb.Config, error) {
	dat, err := os.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file(%s): %v", fn, err)
	}
	var cfg pb.Config
	if err := prototext.Unmarshal(dat, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config proto: %v", err)
	}
	return &cfg, nil
}

func main() {
	flag.Parse()
	if *config == "" {
		log.Fatalf("Provide a config path.")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		glog.Fatalf("failed to listen(): %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxMsgSize(maxMsgSize),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	config, err := readConfig(*config)
	server := new(*store, *port, config)
	fmt.Printf("Will listen on port: %d\n", server.config.GetPort())
	for _, client := range server.config.GetClient() {
		fmt.Printf("Id: %s store: %s\n", client.GetId(), client.GetStore())
	}

	pb.RegisterPiCamServer(s, server)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen&serve: %v", err)
	}

}
