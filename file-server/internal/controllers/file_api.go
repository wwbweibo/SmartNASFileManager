package controllers

import "github.com/gin-gonic/gin"

type FileApiControllers struct {
}

func NewFileApiControllers() *FileApiControllers {
	return &FileApiControllers{}
}

func (f *FileApiControllers) InitRoute(r *gin.Engine) {
	r.GET("/api/v1/file")
}


func (f *FileApiControllers) listDir(c *gin.Context) {

}
