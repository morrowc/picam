module server

go 1.18

replace github.com/morrowc/picam/proto/picam => ./proto

replace github.com/morrowc/picam/client/client => ./client

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/morrowc/picam/client/client v0.0.0-00010101000000-000000000000
	github.com/morrowc/picam/proto/picam v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3 // indirect
)
