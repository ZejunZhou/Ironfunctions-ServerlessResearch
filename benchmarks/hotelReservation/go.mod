module github.com/harlow/go-micro-services

go 1.14

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/protobuf v1.5.2
	github.com/hailocab/go-geoindex v0.0.0-20160127134810-64631bfe9711
	github.com/stretchr/testify v1.9.0 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

// replace cs.utexas.edu/zjia/faas => /src/nightcore/worker/golang

// replace cs.utexas.edu/zjia/faas => ./worker/golang
