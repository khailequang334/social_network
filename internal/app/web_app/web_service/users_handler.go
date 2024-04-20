package web_service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (svc *WebService) CreateUser(ctx *gin.Context) {
	var request types.CreateUserRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
		return
	}

	dob, err := time.Parse("2006-01-02", request.Dob)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
		return
	}

	response, err := svc.UserAndPostClient.CreateUser(ctx, &user_and_post.UserDetailInfo{
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Dob:          timestamppb.New(dob),
		UserName:     request.UserName,
		UserPassword: request.Password,
		Email:        request.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &types.MessageResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &types.MessageResponse{Message: fmt.Sprintf("Successfully created user with id: %d", response.Info.UserId)})
}

func (svc *WebService) EditUser(ctx *gin.Context) {
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

	var request types.EditUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
		return
	}

	// only update requested fields, will be nullptr if not requested
	editUserRequest := &user_and_post.EditUserRequest{}
	editUserRequest.UserId = currentUserId
	if request.FirstName != nil {
		editUserRequest.FirstName = request.FirstName
	}
	if request.LastName != nil {
		editUserRequest.LastName = request.LastName
	}
	if request.Dob != nil {
		parsedDob, err := time.Parse("2006-01-02", *request.Dob)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, &types.MessageResponse{Message: err.Error()})
			return
		}
		editUserRequest.Dob = timestamppb.New(parsedDob)
	}
	if request.Password != nil {
		editUserRequest.UserPassword = request.Password
	}

	// request to user_and_post svc
	response, err := svc.UserAndPostClient.EditUser(ctx, editUserRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &types.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &types.MessageResponse{Message: "Successfully edited user with id: " + fmt.Sprintf("%d", response.UserId)})
}

func (svc *WebService) AuthentcateUser(ctx *gin.Context) {
	start := time.Now()
	countExporter.WithLabelValues("authenticate_user", "total").Inc()

	var status = http.StatusOK
	var request types.LoginRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		countExporter.WithLabelValues("authenticate_user", "bad_request").Inc()
		status = http.StatusBadRequest
		ctx.JSON(status, &types.MessageResponse{Message: err.Error()})
		return
	}

	// request to user_and_post svc
	response, err := svc.UserAndPostClient.AuthenticateUser(ctx, &user_and_post.AuthenticateUserRequest{
		UserName:     request.UserName,
		UserPassword: request.Password,
	})
	if err != nil {
		countExporter.WithLabelValues("authenticate_user", "call_api_failed").Inc()
		status = http.StatusInternalServerError
		ctx.JSON(status, &types.MessageResponse{Message: err.Error()})
		return
	}
	if response.GetStatus() == user_and_post.AuthenticateUserResponse_OK {
		countExporter.WithLabelValues("authenticate_user", "success").Inc()
		ctx.JSON(http.StatusOK, &types.MessageResponse{Message: "ok"})

		// use user id as session id
		ctx.SetCookie("session_id", fmt.Sprintf("%d", response.UserId), 0, "", "", false, false)
	} else if response.GetStatus() == user_and_post.AuthenticateUserResponse_USER_NOT_FOUND {
		countExporter.WithLabelValues("authenticate_user", "not_found").Inc()
		ctx.JSON(http.StatusOK, &types.MessageResponse{Message: "not found"})
	} else {
		countExporter.WithLabelValues("authenticate_user", "wrong_password").Inc()
		ctx.JSON(http.StatusOK, &types.MessageResponse{Message: "wrong password"})
	}

	// defer call when function return
	defer func() {
		latencyExporter.WithLabelValues("authenticate_user", strconv.Itoa(http.StatusOK)).Observe(float64(start.UnixMilli()))
	}()
}
