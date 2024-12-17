package entity

import "fileserver/internal/domain/file"

type ImageCompressionTask struct {
	File file.File
}

func (t *ImageCompressionTask) GetTaskName() string {
	return "ImageCompressionTask"
}
