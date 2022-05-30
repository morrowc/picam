package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

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
	config = flag.String("config", "", "Configuration for the server.")
)

// server holds connection/client information necessary to operate a gRPC server.
type server struct {
	port   int32             // port upon which to listen.
	config *pb.Config        // configuration content for server.
	stores map[string]string // id -> store.
	pb.UnimplementedPiCamServer
	mu sync.Mutex
}

func new(config *pb.Config) *server {
	return &server{
		port:   config.GetPort(),
		config: config,
		stores: make(map[string]string),
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
		glog.Fatalf("Provide a config path.")
	}
	config, err := readConfig(*config)
	server := new(config)

	// Create the port listener.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetPort()))
	if err != nil {
		glog.Fatalf("failed to listen(): %v", err)
	}

	// Set some basic gRPC server options (file size for snd/recv).
	s := grpc.NewServer(
		grpc.MaxMsgSize(maxMsgSize),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	fmt.Printf("Will listen on port: %d\n", server.config.GetPort())
	for _, client := range server.config.GetClient() {
		id := client.GetId()
		dir := client.GetStore()
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			glog.Fatalf("store: %s does not exist: %v", dir, err)
		}
		fmt.Printf("Id: %s store: %s\n", client.GetId(), client.GetStore())
		server.mu.Lock()
		server.stores[id] = dir
		server.mu.Unlock()
	}

	pb.RegisterPiCamServer(s, server)
	reflection.Register(s)

	glog.Info("Ready to start serving.")
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to listen&serve: %v", err)
	}

}
