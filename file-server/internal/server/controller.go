package server

import "github.com/gin-gonic/gin"

type Controller interface {
	InitRoute(r *gin.Engine)
}
