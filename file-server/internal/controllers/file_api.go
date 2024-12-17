package controllers

import (
	"fileserver/internal/biz"
	viewmodel "fileserver/internal/controllers/view_model"

	"github.com/gin-gonic/gin"
)

type FileApiControllers struct {
	fileService *biz.FilerService
}

func NewFileApiControllers(fileService *biz.FilerService) *FileApiControllers {
	return &FileApiControllers{
		fileService: fileService,
	}
}

func (f *FileApiControllers) InitRoute(r *gin.Engine) {
	r.GET("/api/v1/dir", f.dirTree)
	r.GET("/api/v1/file", f.listDir)
	r.POST("/api/v1/file", f.listDir)
	r.GET("/api/v1/file/group", f.listFileByGroup)
}

func (f *FileApiControllers) dirTree(c *gin.Context) {
	resp, err := f.fileService.ListDirectoryTree(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (f *FileApiControllers) listDir(c *gin.Context) {
	req := viewmodel.ListFileRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := f.fileService.ListFiles(c, req.Path)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (f *FileApiControllers) listFileByGroup(ctx *gin.Context) {

}
