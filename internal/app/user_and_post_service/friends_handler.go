package user_and_post_service

import (
	"context"
	"errors"

	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (uaps *UserAndPostService) ensureUserExist(userId int64) error {
	var user types.User
	err := uaps.DB.Table("user").Where("id = ?", userId).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User not found
		uaps.Logger.Error("User not found", zap.Int64("userID", userId))
		return errors.New("user not found")
	} else if err != nil {
		// Other unexpected error
		uaps.Logger.Error("Unknown error while querying user", zap.Error(err))
		return err
	}
	return nil
}

func (uaps *UserAndPostService) FollowUser(ctx context.Context, request *user_and_post.FollowUserRequest) (*user_and_post.FollowUserResponse, error) {
	// Ensure the user exists
	err := uaps.ensureUserExist(request.UserId)
	if err != nil {
		return &user_and_post.FollowUserResponse{Status: user_and_post.FollowUserResponse_USER_NOT_FOUND}, nil
	}

	// Ensure the following user exists
	err = uaps.ensureUserExist(request.FollowingUserId)
	if err != nil {
		return &user_and_post.FollowUserResponse{Status: user_and_post.FollowUserResponse_USER_NOT_FOUND}, nil
	}

	// Preload user with their following and follower relationships
	var user types.User
	err = uaps.DB.Preload("Following").Preload("Follower").First(&user, request.UserId).Error
	if err != nil {
		return nil, errors.New("error fetching user")
	}
	uaps.Logger.Debug("Returned user", zap.Any("user", user))

	// Check if the user is already being followed
	var alreadyFollowed bool
	for _, following := range user.Following {
		if following.ID == uint(request.FollowingUserId) {
			alreadyFollowed = true
			break
		}
	}
	if alreadyFollowed {
		return &user_and_post.FollowUserResponse{Status: user_and_post.FollowUserResponse_ALREADY_FOLLOWED}, nil
	}

	// Fetch the user to follow
	var followingUser types.User
	err = uaps.DB.First(&followingUser, request.FollowingUserId).Error
	if err != nil {
		return nil, errors.New("error fetching following user")
	}

	// Append following user to user's following relationship
	err = uaps.DB.Model(&user).Association("Following").Append(&followingUser)
	if err != nil {
		return nil, errors.New("error appending following user")
	}

	// Append user to following user's follower relationship
	err = uaps.DB.Model(&followingUser).Association("Follower").Append(&user)
	if err != nil {
		return nil, errors.New("error appending follower user")
	}

	uaps.Logger.Info("following new user")
	return &user_and_post.FollowUserResponse{Status: user_and_post.FollowUserResponse_OK}, nil

}

func (uaps *UserAndPostService) UnfollowUser(ctx context.Context, request *user_and_post.UnfollowUserRequest) (*user_and_post.UnfollowUserResponse, error) {
	err := uaps.ensureUserExist(request.UserId)
	if err != nil {
		return &user_and_post.UnfollowUserResponse{Status: user_and_post.UnfollowUserResponse_USER_NOT_FOUND}, nil
	}
	err = uaps.ensureUserExist(request.FollowingUserId)
	if err != nil {
		return &user_and_post.UnfollowUserResponse{Status: user_and_post.UnfollowUserResponse_USER_NOT_FOUND}, nil
	}

	var user types.User
	uaps.DB.Preload("Following").Preload("Follower").Find(&user, request.UserId)
	uaps.Logger.Debug("returned user", zap.Any("user", user))

	var alreadyFollowed bool
	for _, following := range user.Following {
		if following.ID == uint(request.FollowingUserId) {
			alreadyFollowed = true
			break
		}
	}

	if alreadyFollowed {
		uaps.Logger.Info("unfollowing user")
		var followingUser types.User
		uaps.DB.Where(&types.User{ID: uint(request.FollowingUserId)}).First(&followingUser)

		err := uaps.DB.Model(&user).Association("Following").Delete(&followingUser)
		if err != nil {
			return nil, err
		}
		err = uaps.DB.Model(&followingUser).Association("Follower").Delete(&user)
		if err != nil {
			return nil, err
		}
		return &user_and_post.UnfollowUserResponse{Status: user_and_post.UnfollowUserResponse_OK}, nil
	} else {
		return &user_and_post.UnfollowUserResponse{Status: user_and_post.UnfollowUserResponse_NOT_FOLLOWED}, nil
	}
}

func (uaps *UserAndPostService) GetFollowerList(ctx context.Context, req *user_and_post.GetFollowerListRequest) (*user_and_post.GetFollowerListResponse, error) {
	err := uaps.ensureUserExist(req.UserId)
	if err != nil {
		return &user_and_post.GetFollowerListResponse{Status: user_and_post.GetFollowerListResponse_USER_NOT_FOUND}, nil
	}

	var user types.User
	uaps.DB.Preload("Follower").Find(&user, req.UserId)
	uaps.Logger.Debug("returned user", zap.Any("user", user))
	var followersInfo []*user_and_post.GetFollowerListResponse_FollowerInfo
	for _, follower := range user.Follower {
		followersInfo = append(followersInfo, &user_and_post.GetFollowerListResponse_FollowerInfo{
			UserId:   int64(follower.ID),
			UserName: follower.UserName,
		})
	}
	return &user_and_post.GetFollowerListResponse{Status: user_and_post.GetFollowerListResponse_OK, Followers: followersInfo}, nil
}
