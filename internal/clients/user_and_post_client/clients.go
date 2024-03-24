package user_and_post_client

import (
	"context"
	"math/rand"

	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type randomClients struct {
	clients []user_and_post.UserAndPostClient
}

func (a *randomClients) CreateUser(ctx context.Context, in *user_and_post.UserDetailInfo, opts ...grpc.CallOption) (*user_and_post.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CreateUser(ctx, in, opts...)
}

func (a *randomClients) EditUser(ctx context.Context, in *user_and_post.EditUserRequest, opts ...grpc.CallOption) (*user_and_post.EditUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].EditUser(ctx, in, opts...)
}

func (a *randomClients) AuthenticateUser(ctx context.Context, in *user_and_post.AuthenticateUserRequest, opts ...grpc.CallOption) (*user_and_post.AuthenticateUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].AuthenticateUser(ctx, in, opts...)
}

func (a *randomClients) FollowUser(ctx context.Context, in *user_and_post.FollowUserRequest, opts ...grpc.CallOption) (*user_and_post.FollowUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].FollowUser(ctx, in, opts...)
}

func (a *randomClients) UnfollowUser(ctx context.Context, in *user_and_post.UnfollowUserRequest, opts ...grpc.CallOption) (*user_and_post.UnfollowUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].UnfollowUser(ctx, in, opts...)
}

func (a *randomClients) GetFollowerList(ctx context.Context, in *user_and_post.GetFollowerListRequest, opts ...grpc.CallOption) (*user_and_post.GetFollowerListResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetFollowerList(ctx, in, opts...)
}

func (a *randomClients) CreatePost(ctx context.Context, in *user_and_post.CreatePostRequest, opts ...grpc.CallOption) (*user_and_post.CreatePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CreatePost(ctx, in, opts...)
}

func (a *randomClients) GetPost(ctx context.Context, in *user_and_post.GetPostRequest, opts ...grpc.CallOption) (*user_and_post.GetPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetPost(ctx, in, opts...)
}

func (a *randomClients) DeletePost(ctx context.Context, in *user_and_post.DeletePostRequest, opts ...grpc.CallOption) (*user_and_post.DeletePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].DeletePost(ctx, in, opts...)
}

func (a *randomClients) EditPost(ctx context.Context, in *user_and_post.EditPostRequest, opts ...grpc.CallOption) (*user_and_post.EditPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].EditPost(ctx, in, opts...)
}

func (a *randomClients) LikePost(ctx context.Context, in *user_and_post.LikePostRequest, opts ...grpc.CallOption) (*user_and_post.LikePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].LikePost(ctx, in, opts...)
}

func (a *randomClients) CommentPost(ctx context.Context, in *user_and_post.CommentPostRequest, opts ...grpc.CallOption) (*user_and_post.CommentPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CommentPost(ctx, in, opts...)
}

func NewClients(hosts []string) (user_and_post.UserAndPostClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]user_and_post.UserAndPostClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			return nil, err
		}
		client := user_and_post.NewUserAndPostClient(conn)
		clients = append(clients, client)
	}
	return &randomClients{clients}, nil
}
