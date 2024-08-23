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
