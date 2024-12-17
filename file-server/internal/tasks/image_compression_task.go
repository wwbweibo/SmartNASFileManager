package tasks

import (
	"context"
	"fileserver/internal/domain/file"
	"fileserver/internal/server"
	"fileserver/internal/tasks/entity"
	"fileserver/utils"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

type ImageCompressionTaskHandler struct {
	imageChan   chan file.File
	cacheDir    string
	startAt     time.Time
	nasRootPath string
}

func NewImageCompressionTaskHandler(nasRootPath, cachePath string) *ImageCompressionTaskHandler {
	handler := &ImageCompressionTaskHandler{
		imageChan:   make(chan file.File),
		cacheDir:    cachePath,
		nasRootPath: nasRootPath,
	}
	bus.RegisterHandler(handler)
	return handler
}

func (t *ImageCompressionTaskHandler) GetTaskName() string {
	return "ImageCompressionTask"
}

func (t *ImageCompressionTaskHandler) GetRunningDuration() time.Duration {
	return time.Since(t.startAt)
}

func (t *ImageCompressionTaskHandler) Start(ctx context.Context) error {
	t.startAt = time.Now()
	for {
		select {
		case image := <-t.imageChan:
			// compress image
			t.compressImage(image)
		case <-ctx.Done():
			return nil
		}
	}
}

func (t *ImageCompressionTaskHandler) Append(task server.ITask) {
	if task.GetTaskName() != t.GetTaskName() {
		return
	}
	tt := task.(*entity.ImageCompressionTask)
	t.imageChan <- tt.File
}

func (t *ImageCompressionTaskHandler) Stop(ctx context.Context) error {
	close(t.imageChan)
	return nil
}

func (t *ImageCompressionTaskHandler) compressImage(_image file.File) {
	// compress image
	file, err := os.Open(t.nasRootPath + _image.Path)
	if err != nil {
		log.Printf("error opening file %s: %v", _image.Path, err)
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("error decoding image %s: %v", _image.Path, err)
		return
	}
	resized := resize.Thumbnail(512, 512, img, resize.Lanczos3)
	savePath := t.cacheDir + _image.Path
	saveDir := utils.GetDirectory(savePath)
	_, err = os.Stat(saveDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(saveDir, 0755)
			if err != nil {
				log.Printf("error creating directory %s: %v", saveDir, err)
				return
			}
		} else {
			log.Printf("error checking directory %s: %v", saveDir, err)
			return
		}
	}
	out, err := os.Create(savePath)
	if err != nil {
		log.Printf("error creating file %s: %v", t.cacheDir+_image.Path, err)
		return
	}
	defer out.Close()
	err = jpeg.Encode(out, resized, nil)
	if err != nil {
		log.Printf("error encoding image %s: %v", _image.Path, err)
	}
}
