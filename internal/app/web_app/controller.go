package web_app

import (
	"fmt"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/khailequang334/social_network/internal/app/web_app/web_service"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type WebController struct {
	web_service.WebService
	Port int
}

func (c WebController) Run() {
	r := gin.Default()

	// add router
	v1 := r.Group("/api/v1")
	AddRouter(v1, &c.WebService)

	AddSwagger(r)
	AddPrometheus(r)
	AddPprof(r)

	r.Run(fmt.Sprintf(":%d", c.Port))
}

func AddRouter(r *gin.RouterGroup, svc *web_service.WebService) {
	// router for users handler
	userRouter := r.Group("users")
	userRouter.POST("", svc.CreateUser)
	userRouter.POST("login", svc.AuthentcateUser)
	userRouter.PUT("", svc.EditUser)

	// router for friends handler
	friendRouter := r.Group("friends")
	friendRouter.GET(":user_id", svc.GetFollowList)
	friendRouter.POST(":user_id", svc.FollowUser)
	friendRouter.DELETE(":user_id", svc.UnfollowUser)

	// router for posts handler
	postRouter := r.Group("posts")
	postRouter.POST("", svc.CreatePost)
	postRouter.GET(":post_id", svc.GetPost)
	postRouter.PUT(":post_id", svc.EditPost)
	postRouter.DELETE(":post_id", svc.DeletePost)
	postRouter.POST(":post_id/likes", svc.LikePost)
	postRouter.POST(":post_id/comments", svc.CreatePostComment)

	// router for newsfeeds handler
	newsfeedRouter := r.Group("newsfeeds")
	newsfeedRouter.GET("", svc.GetNewsfeed)
}

func AddSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func AddPrometheus(r *gin.Engine) {
	handler := promhttp.Handler()
	r.GET("/metrics", func(context *gin.Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	})
}

func AddPprof(r *gin.Engine) {
	r.GET("/debug/pprof/", func(context *gin.Context) {
		pprof.Index(context.Writer, context.Request)
	})

	r.GET("/debug/pprof/profile", func(context *gin.Context) {
		pprof.Profile(context.Writer, context.Request)
	})

	r.GET("/debug/pprof/trace", func(context *gin.Context) {
		pprof.Trace(context.Writer, context.Request)
	})
}
