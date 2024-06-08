package web_server

import (
	"fmt"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/khailequang334/social_network/internal/interfaces/app/web_server/web_service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type WebServer struct {
	Service *web_service.WebService
	Port    int
}

func (s *WebServer) Run() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	setupRoutes(v1, s.Service)

	setupPrometheus(r)
	setupPprof(r)

	r.Run(fmt.Sprintf(":%d", s.Port))
}

func setupRoutes(r *gin.RouterGroup, svc *web_service.WebService) {
	userRouter := r.Group("users")
	userRouter.POST("", svc.CreateUser)
	userRouter.POST("login", svc.AuthentcateUser)
	userRouter.PUT("", svc.EditUser)

	friendRouter := r.Group("friends")
	friendRouter.GET(":user_id", svc.GetFollowList)
	friendRouter.POST(":user_id", svc.FollowUser)
	friendRouter.DELETE(":user_id", svc.UnfollowUser)

	postRouter := r.Group("posts")
	postRouter.POST("", svc.CreatePost)
	postRouter.GET(":post_id", svc.GetPost)
	postRouter.PUT(":post_id", svc.EditPost)
	postRouter.DELETE(":post_id", svc.DeletePost)
	postRouter.POST(":post_id/likes", svc.LikePost)
	postRouter.POST(":post_id/comments", svc.CreatePostComment)

	newsfeedRouter := r.Group("newsfeeds")
	newsfeedRouter.GET("", svc.GetNewsfeed)
}

func setupPrometheus(r *gin.Engine) {
	handler := promhttp.Handler()
	r.GET("/metrics", func(context *gin.Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	})
}

func setupPprof(r *gin.Engine) {
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
