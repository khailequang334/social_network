package web_service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khailequang334/social_network/internal/interfaces/proto/protobuf/newsfeed"
	"github.com/khailequang334/social_network/internal/model"
	"go.uber.org/zap"
)

func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		svc.Logger.Error("invalid user id", zap.String("user_id", ctx.Param("user_id")))
		ctx.JSON(http.StatusBadRequest, model.MessageResponse{Message: "invalid user id"})
		return
	}

	response, err := svc.NewsfeedClient.GenerateNewsfeed(ctx, &newsfeed.GenerateNewsfeedRequest{
		UserId: userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.GetPostIds())
}
