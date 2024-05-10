package user_and_post_service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/types"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Retrieve post from Redis cache
func (uaps *UserAndPostService) getPostFromCache(ctx context.Context, postId int64) (*user_and_post.GetPostResponse, error) {
	cacheKey := "post:" + strconv.FormatInt(postId, 10)
	cachedData, err := uaps.Redis.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cachedPost user_and_post.GetPostResponse
	err = json.Unmarshal(cachedData, &cachedPost)
	if err != nil {
		return nil, err
	}

	return &cachedPost, nil
}

// Cache post in Redis
func (uaps *UserAndPostService) cachePost(ctx context.Context, postId int64, postData *user_and_post.GetPostResponse) error {
	cacheKey := "post:" + strconv.FormatInt(postId, 10)
	cacheDuration := time.Hour // cache duration 1 hour
	data, err := json.Marshal(postData)
	if err != nil {
		return err
	}

	return uaps.Redis.Set(ctx, cacheKey, data, cacheDuration).Err()
}

func (uaps *UserAndPostService) GetPost(ctx context.Context, request *user_and_post.GetPostRequest) (*user_and_post.GetPostResponse, error) {
	uaps.Logger.Debug("start get post")
	defer uaps.Logger.Debug("end get post")

	// Check if the post exists in Redis cache
	cachedPost, err := uaps.getPostFromCache(ctx, request.PostId)
	if err != nil {
		uaps.Logger.Error("failed to get post from cache", zap.Error(err), zap.Int64("PostId", request.PostId))
	}
	if cachedPost != nil {
		return cachedPost, nil
	}

	var post types.Post
	err = uaps.DB.First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user_and_post.GetPostResponse{
			Status: user_and_post.GetPostResponse_POST_NOT_FOUND,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	// Cache the retrieved post
	err = uaps.cachePost(ctx, request.PostId, &user_and_post.GetPostResponse{
		Status: user_and_post.GetPostResponse_OK,
		Post: &user_and_post.Post{
			PostId:           int64(post.ID),
			UserId:           int64(post.UserID),
			ContentText:      post.ContentText,
			ContentImagePath: post.ContentImagePath,
			Visible:          post.Visible,
			CreatedTime:      timestamppb.New(post.CreatedAt),
		},
	})
	if err != nil {
		uaps.Logger.Error("failed to cache post", zap.Error(err), zap.Int64("PostId", request.PostId))
	}

	return &user_and_post.GetPostResponse{
		Status: user_and_post.GetPostResponse_OK,
		Post: &user_and_post.Post{
			PostId:           int64(post.ID),
			UserId:           int64(post.UserID),
			ContentText:      post.ContentText,
			ContentImagePath: post.ContentImagePath,
			Visible:          post.Visible,
			CreatedTime:      timestamppb.New(post.CreatedAt),
		},
	}, nil
}

func (uaps *UserAndPostService) CreatePost(ctx context.Context, request *user_and_post.CreatePostRequest) (*user_and_post.CreatePostResponse, error) {
	uaps.Logger.Debug("start create post")
	defer uaps.Logger.Debug("end create post")

	err := uaps.ensureUserExist(request.UserId)
	if err != nil {
		return &user_and_post.CreatePostResponse{
			Status: user_and_post.CreatePostResponse_USER_NOT_FOUND,
		}, nil
	}

	var user types.User
	err = uaps.DB.Preload("Posts").Find(&user, request.UserId).Error
	if err != nil {
		return nil, err
	}

	post := types.Post{
		ContentText:      request.ContentText,
		ContentImagePath: request.ContentImagePath,
		UserID:           uint(request.UserId),
		Visible:          request.Visible,
	}

	err = uaps.DB.Model(&user).Association("Posts").Append(&post)
	if err != nil {
		return nil, err
	}
	return &user_and_post.CreatePostResponse{Status: user_and_post.CreatePostResponse_OK, PostId: int64(post.ID)}, nil
}

func (uaps *UserAndPostService) DeletePost(ctx context.Context, request *user_and_post.DeletePostRequest) (*user_and_post.DeletePostResponse, error) {
	uaps.Logger.Debug("start delete post")
	defer uaps.Logger.Debug("end delete post")

	var post types.Post
	err := uaps.DB.First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user_and_post.DeletePostResponse{
			Status: user_and_post.DeletePostResponse_POST_NOT_FOUND,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	err = uaps.DB.Delete(&post).Error
	if err != nil {
		return nil, err
	}
	return &user_and_post.DeletePostResponse{Status: user_and_post.DeletePostResponse_OK}, nil
}

// EditPost edit post
func (uaps *UserAndPostService) EditPost(ctx context.Context, request *user_and_post.EditPostRequest) (*user_and_post.EditPostResponse, error) {
	uaps.Logger.Debug("start edit post", zap.Any("request", request))
	defer uaps.Logger.Debug("end edit post")

	var post types.Post
	err := uaps.DB.First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user_and_post.EditPostResponse{
			Status: user_and_post.EditPostResponse_POST_NOT_FOUND,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	uaps.Logger.Debug("post", zap.Any("post", post))

	if request.ContentText != nil {
		post.ContentText = request.GetContentText()
	}
	if request.ContentImagePath != nil {
		post.ContentImagePath = request.GetContentImagePath()
	}
	if request.Visible != nil {
		post.Visible = request.GetVisible()
	}
	uaps.Logger.Debug("updated post", zap.Any("post", post))

	err = uaps.DB.Save(&post).Error
	if err != nil {
		return nil, err
	}
	return &user_and_post.EditPostResponse{Status: user_and_post.EditPostResponse_OK}, nil
}

// CreatePostComment create post comment
func (uaps *UserAndPostService) CreatePostComment(ctx context.Context, request *user_and_post.CommentPostRequest) (*user_and_post.CommentPostResponse, error) {
	uaps.Logger.Debug("start create post comment")
	defer uaps.Logger.Debug("end create post comment")

	err := uaps.ensureUserExist(request.UserId)
	if err != nil {
		return &user_and_post.CommentPostResponse{
			Status: user_and_post.CommentPostResponse_USER_NOT_FOUND,
		}, nil
	}

	var post types.Post
	err = uaps.DB.Where("id = ?", request.PostId).First(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		uaps.Logger.Debug("post not found")
		return &user_and_post.CommentPostResponse{
			Status: user_and_post.CommentPostResponse_POST_NOT_FOUND,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	uaps.Logger.Debug("post found", zap.Any("post", post))
	uaps.DB.Preload("Comments").Find(&post, request.PostId)

	comment := types.Comment{
		Content: request.ContentText,
		UserID:  uint(request.UserId),
		PostID:  uint(request.PostId),
	}

	err = uaps.DB.Model(&post).Association("Comments").Append(&comment)
	if err != nil {
		return nil, err
	}
	return &user_and_post.CommentPostResponse{Status: user_and_post.CommentPostResponse_OK, CommentId: int64(comment.ID)}, nil
}

// LikePost like post
func (uaps *UserAndPostService) LikePost(ctx context.Context, request *user_and_post.LikePostRequest) (*user_and_post.LikePostResponse, error) {
	uaps.Logger.Debug("start like post")
	defer uaps.Logger.Debug("end like post")

	var err error
	var user types.User
	err = uaps.DB.First(&user, request.UserId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user_and_post.LikePostResponse{
			Status: user_and_post.LikePostResponse_USER_NOT_FOUND,
		}, nil
	}

	var post types.Post
	err = uaps.DB.First(&post, request.PostId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user_and_post.LikePostResponse{
			Status: user_and_post.LikePostResponse_POST_NOT_FOUND,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	err = uaps.DB.Preload("LikedUsers").Find(&post, request.PostId).Error
	if err != nil {
		return nil, err
	}

	err = uaps.DB.Model(&post).Association("LikedUsers").Append(&user)
	if err != nil {
		return nil, err
	}
	return &user_and_post.LikePostResponse{Status: user_and_post.LikePostResponse_OK}, nil
}
