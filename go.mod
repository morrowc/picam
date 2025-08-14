module server

go 1.23.0

replace github.com/morrowc/picam/proto/picam => ./proto

replace github.com/morrowc/picam/client/client => ./client

require (
	github.com/golang/glog v1.2.4
	github.com/morrowc/picam/client/client v0.0.0-00010101000000-000000000000
	github.com/morrowc/picam/proto/picam v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
)
