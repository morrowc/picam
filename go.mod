module server

go 1.18

replace github.com/morrowc/picam/proto/picam => ./proto

replace github.com/morrowc/picam/client/client => ./client

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/morrowc/picam/client/client v0.0.0-00010101000000-000000000000
)

require (
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gidoBOSSftw5731/log v0.0.0-20210527210830-1611311b4b64 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/morrowc/picam/proto/picam v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-20220516155154-20f960328961 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3 // indirect
	google.golang.org/grpc v1.46.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
