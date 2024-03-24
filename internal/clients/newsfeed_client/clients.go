package newsfeed_client

import (
	"context"
	"math/rand"

	"github.com/khailequang334/social_network/internal/protobuf/newsfeed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type randomClients struct {
	clients []newsfeed.NewsfeedClient
}

func (r *randomClients) GenerateNewsfeed(ctx context.Context, in *newsfeed.NewsfeedRequest, opts ...grpc.CallOption) (*newsfeed.NewsfeedResponse, error) {
	return r.clients[rand.Intn(len(r.clients))].GenerateNewsfeed(ctx, in, opts...)
}

func NewClients(hosts []string) (newsfeed.NewsfeedClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]newsfeed.NewsfeedClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			return nil, err
		}
		client := newsfeed.NewNewsfeedClient(conn)
		clients = append(clients, client)
	}
	return &randomClients{clients}, nil
}
