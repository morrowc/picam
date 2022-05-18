module server

go 1.18

require (
	github.com/gidoBOSSftw5731/log v0.0.0-20210527210830-1611311b4b64
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/morrowc/picam/proto/picam v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/morrowc/picam/proto/picam => ./proto
