package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	engine      *gin.Engine
	controllers []Controller
	server      *http.Server
}

func NewGinServer() *GinServer {
	r := gin.Default()
	// 添加跨域
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Next()
	})
	return &GinServer{
		engine: r,
	}
}

func (g *GinServer) Start(ctx context.Context) error {
	// start a http server here
	for _, controller := range g.controllers {
		controller.InitRoute(g.engine)
	}
	go func() {
		<-ctx.Done()
		_ = g.server.Shutdown(ctx)
	}()
	g.server = &http.Server{
		Addr:    ":8080",
		Handler: g.engine.Handler(),
	}
	return g.server.ListenAndServe()
}

func (g *GinServer) RegisterController(controllers ...Controller) {
	g.controllers = append(g.controllers, controllers...)
}
