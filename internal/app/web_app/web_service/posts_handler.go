package web_service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/types"
)

func (svc *WebService) CreatePost(ctx *gin.Context) {
	var request types.CreatePostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	response, err := svc.UserAndPostClient.CreatePost(ctx, &user_and_post.CreatePostRequest{
		UserId:           request.UserId,
		ContentText:      request.ContentText,
		ContentImagePath: request.ContentImagePath,
		Visible:          request.Visible,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.CreatePostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: fmt.Sprintf("create post successfully with id: %d", response.PostId)})
}

func (svc *WebService) GetPost(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: fmt.Sprintf("invalid post id: %s", ctx.Param("post_id"))})
		return
	}

	response, err := svc.UserAndPostClient.GetPost(ctx, &user_and_post.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.GetPostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, types.PostDetailResponse{
		PostID:           response.Post.PostId,
		UserID:           response.Post.UserId,
		ContentText:      response.Post.ContentText,
		ContentImagePath: response.Post.ContentImagePath,
		Visible:          response.Post.Visible,
		CreatedTime:      response.Post.CreatedTime.AsTime(),
	})
}

func (svc *WebService) DeletePost(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: fmt.Sprintf("invalid post id: %s", ctx.Param("post_id"))})
		return
	}

	response, err := svc.UserAndPostClient.DeletePost(ctx, &user_and_post.DeletePostRequest{
		PostId: postId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.DeletePostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: fmt.Sprintf("delete post successfully with id: %d", postId)})
}

func (svc *WebService) EditPost(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: fmt.Sprintf("invalid post id: %s", ctx.Param("post_id"))})
		return
	}

	var request types.EditPostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	editPostRequest := &user_and_post.EditPostRequest{}
	editPostRequest.PostId = postId
	if request.ContentText != nil {
		editPostRequest.ContentText = request.ContentText
	}
	if request.ContentImagePath != nil {
		editPostRequest.ContentImagePath = request.ContentImagePath
	}
	if request.Visible != nil {
		editPostRequest.Visible = request.Visible
	}

	response, err := svc.UserAndPostClient.EditPost(ctx, editPostRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.EditPostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: fmt.Sprintf("edit post successfully with id: %d", postId)})
}

func (svc *WebService) LikePost(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: fmt.Sprintf("invalid post id: %s", ctx.Param("post_id"))})
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

	response, err := svc.UserAndPostClient.LikePost(ctx, &user_and_post.LikePostRequest{
		PostId: postId,
		UserId: currentUserId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if response.Status == user_and_post.LikePostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if response.Status == user_and_post.LikePostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: fmt.Sprintf("like post successfully with id: %d", postId)})
}

func (svc *WebService) CreatePostComment(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: fmt.Sprintf("invalid post id: %s", ctx.Param("post_id"))})
		return
	}

	var request types.CreatePostCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
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

	resp, err := svc.UserAndPostClient.CommentPost(ctx, &user_and_post.CommentPostRequest{
		PostId:      postId,
		UserId:      currentUserId,
		ContentText: request.ContentText,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.Status == user_and_post.CommentPostResponse_POST_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.Status == user_and_post.CommentPostResponse_USER_NOT_FOUND {
		ctx.JSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, types.MessageResponse{Message: fmt.Sprintf("create post comment successfully with id: %d", resp.CommentId)})
}
