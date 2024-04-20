package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/app/user_and_post_service"
	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"google.golang.org/grpc"
)

var (
	path = flag.String("conf", "config.yml", "config path for this service")
)

func main() {
	// Start authenticate and post service
	conf, err := configs.GetUserAndPostConfig(*path)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	service, err := user_and_post_service.NewUserAndPostService(conf)
	if err != nil {
		log.Fatalf("failed to init server %s", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", conf.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	user_and_post.RegisterUserAndPostServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped %v", err)
	}
}
