module server

go 1.18

require (
	github.com/gidoBOSSftw5731/log v0.0.0-20210527210830-1611311b4b64
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/morrowc/picam/proto/picam v0.0.0-00010101000000-000000000000
)

replace github.com/morrowc/picam/proto/picam => ./proto
