package web_service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (svc *WebService) FollowUser(ctx *gin.Context) {
	followingUserId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		svc.Logger.Error("invalid user id", zap.String("user_id", ctx.Param("user_id")))
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "invalid user id"})
		return
	}

	// get current user from session id
	sessionId, err := ctx.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		svc.Logger.Error("unauthorized")
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: "unauthorized"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unexpected error"})
		return
	}

	currentUserId, err := strconv.ParseInt(sessionId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unexpected error"})
		return
	}

	response, err := svc.UserAndPostClient.FollowUser(ctx, &user_and_post.FollowUserRequest{
		UserId:          currentUserId,
		FollowingUserId: followingUserId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.FollowUserResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if response.Status == user_and_post.FollowUserResponse_ALREADY_FOLLOWED {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "user already followed"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: "follow user successfully"})
}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		svc.Logger.Error("invalid user id", zap.String("user_id", ctx.Param("user_id")))
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "invalid user id"})
		return
	}

	// get current user from session id
	sessionId, err := ctx.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		svc.Logger.Error("unauthorized")
		ctx.JSON(http.StatusUnauthorized, types.MessageResponse{Message: "unauthorized"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unexpected error"})
		return
	}

	currentUserId, err := strconv.ParseInt(sessionId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: "unexpected error"})
		return
	}

	response, err := svc.UserAndPostClient.UnfollowUser(ctx, &user_and_post.UnfollowUserRequest{
		UserId:          currentUserId,
		FollowingUserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.UnfollowUserResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if response.Status == user_and_post.UnfollowUserResponse_NOT_FOLLOWED {
		ctx.JSON(http.StatusOK, types.MessageResponse{Message: "user not followed"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: "unfollow user successfully"})
}

func (svc *WebService) GetFollowList(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64) // get int64 value
	if err != nil {
		svc.Logger.Error("invalid user id", zap.String("user_id", ctx.Param("user_id")))
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "invalid user id"})
		return
	}

	response, err := svc.UserAndPostClient.GetFollowerList(ctx, &user_and_post.GetFollowerListRequest{UserId: userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.GetFollowerListResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, response.GetFollowers())
}
