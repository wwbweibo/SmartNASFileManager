package tasks

import (
	"context"
	"fileserver/internal/domain/file"
	"fileserver/utils"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

type ImageCompressionTask struct {
	imageChan   chan file.File
	cacheDir    string
	startAt     time.Time
	nasRootPath string
	cachePath   string
}

func NewImageCompressionTask(nasRootPath, cachePath string) *ImageCompressionTask {
	return &ImageCompressionTask{
		imageChan:   make(chan file.File),
		cacheDir:    cachePath,
		nasRootPath: nasRootPath,
	}
}

func (t *ImageCompressionTask) GetTaskName() string {
	return "ImageCompressionTask"
}

func (t *ImageCompressionTask) GetRunningDuration() time.Duration {
	return time.Since(t.startAt)
}

func (t *ImageCompressionTask) Start(ctx context.Context) error {
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

func (t *ImageCompressionTask) AddImage(image file.File) {
	t.imageChan <- image
}

func (t *ImageCompressionTask) Stop(ctx context.Context) error {
	close(t.imageChan)
	return nil
}

func (t *ImageCompressionTask) compressImage(_image file.File) {
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
	return
}
