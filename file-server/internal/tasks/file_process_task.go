package tasks

import (
	"context"
	"fileserver/internal"
	"fileserver/internal/adapters/dl"
	domainFile "fileserver/internal/domain/file"
	"fileserver/internal/server"
	"fileserver/internal/tasks/entity"
	"fileserver/utils"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// FileProcessTask 文件处理任务
type FileProcessTaskHandler struct {
	startTime time.Time
	option    utils.ScanOptions
	repo      domainFile.IFileRepository
	config    internal.Config
	taskChan  chan server.ITask
}

func NewFileProcessTask(
	option utils.ScanOptions,
	repo domainFile.IFileRepository,
	config internal.Config,
) *FileProcessTaskHandler {
	return &FileProcessTaskHandler{
		option:   option,
		repo:     repo,
		config:   config,
		taskChan: make(chan server.ITask, 100),
	}
}

func (t *FileProcessTaskHandler) GetTaskName() string {
	return "FileProcessTask"
}

func (t *FileProcessTaskHandler) GetRunningDuration() time.Duration {
	return time.Since(t.startTime)
}

func (t *FileProcessTaskHandler) Start(ctx context.Context) error {
	t.startTime = time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		case task := <-t.taskChan:
			_task := (task).(*entity.FileProcessTask)
			t.singleFileHandler(ctx, _task.File)
		}
	}
}

func (t *FileProcessTaskHandler) Stop(ctx context.Context) error {
	return nil
}

func (t *FileProcessTaskHandler) Append(task server.ITask) {
	if task.GetTaskName() != t.GetTaskName() {
		return
	}
	t.taskChan <- task
}

func (t *FileProcessTaskHandler) singleFileHandler(ctx context.Context, file string) {
	log.Default().Printf("handling file %s", file)
	file = strings.Replace(file, t.option.RootPath, "", 1)
	_file := domainFile.NewFile(file)
	_file.Size, _file.UpdatedAt = utils.GetFileSize(t.option.RootPath + file)
	dbFile, _ := t.repo.GetFileByPath(ctx, file)
	if dbFile.UpdatedAt.Equal(_file.UpdatedAt) || dbFile.UpdatedAt.After(_file.UpdatedAt) {
		// 记录的文件更新时间等于或晚于当前文件的更新时间
		// 说明文件没有更新
		return
	}
	// 记录的文件更新时间早于当前文件的更新时间
	_file.Checksum = utils.Sha256(t.option.RootPath + file)
	if dbFile.Checksum == _file.Checksum {
		if !(dbFile.Group == "unknown" || dbFile.Group == "") {
			fmt.Printf("file %s has not changed\n", file)
			return
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := dl.NewClient(t.config.DLConfiguration).Understanding(ctx, dl.UnderstandingRequest{
			Path: file,
		})
		if err != nil {
			log.Default().Printf("error getting file type: %v", err)
			return
		}
		_file.SetFileTypeFromUnderstanding(result)
	}()

	// insert into database
	wg.Wait()
	if _file.Group == "image" {
		bus.Send(&entity.ImageCompressionTask{File: _file})
	}
	err := t.repo.CreateOrUpdateFile(ctx, _file)
	if err != nil {
		log.Default().Printf("error inserting file %s: %v", file, err)
	}
}
