package tasks

import (
	"context"
	domainFile "fileserver/internal/domain/file"
	"fileserver/internal/server"
	"fileserver/internal/tasks/entity"
	"fileserver/utils"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileSystemWatchTask 将在后台监控文件系统的变化，并将变化的文件信息同步到数据库中
type FileSystemWatchTask struct {
	watcher     *fsnotify.Watcher
	opts        utils.ScanOptions
	repo        domainFile.IFileRepository
	nasRootPath string
}

func NewFileSystemWatchTask(nasRootPath string,
	opts utils.ScanOptions,
	repo domainFile.IFileRepository,
) *FileSystemWatchTask {
	return &FileSystemWatchTask{
		nasRootPath: nasRootPath,
		opts:        opts,
		repo:        repo,
	}
}

// implement of BackendTaskHandler interface
func (f *FileSystemWatchTask) GetTaskName() string {
	return "FileSystemWatchTask"
}

func (f *FileSystemWatchTask) GetRunningDuration() time.Duration {
	return 0
}

func (f *FileSystemWatchTask) Start(ctx context.Context) (err error) {
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Default().Printf("failed to create watcher: %s", err)
		return err
	}
	// fsnotify was not watch the sub directory, so, need to walk the directory to find all the sub directory
	_, dirs := utils.WalkDir(f.nasRootPath)
	for _, dir := range dirs {
		log.Default().Printf("add watch path: %s", dir)
		err = f.watcher.Add(dir)
		if err != nil {
			log.Default().Printf("failed to add watch path: %s", err)
			return err
		}
	}
	for {
		select {
		case event := <-f.watcher.Events:
			f.onFileSystemEvent(event)
		case err := <-f.watcher.Errors:
			log.Default().Printf("watcher error: %s", err)
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (f *FileSystemWatchTask) onFileSystemEvent(event fsnotify.Event) {
	log.Default().Printf("file system event: %v", event)
	if !f.opts.ShouldWatch(event.Name, true) {
		return
	}
	switch event.Op {
	case fsnotify.Create:
		f.onFileCreate(event)
	case fsnotify.Remove:
		f.onFileRemove(event)
	case fsnotify.Rename:
		// for rename event, we treat it as remove event, the new file will be created as a new file
		f.onFileRemove(event)
	}

}

func (f *FileSystemWatchTask) onFileRemove(event fsnotify.Event) {
	fileDBName := strings.Replace(event.Name, f.nasRootPath, "", 1)
	f.repo.RemoveFile(context.Background(), fileDBName)
}

func (f *FileSystemWatchTask) onFileCreate(event fsnotify.Event) {
	if utils.CheckIsDir(event.Name) {
		log.Default().Printf("add watch path: %s", event.Name)
		f.watcher.Add(event.Name)
		return
	}
	// TODO: create event may not reliable, since some file will very huge, when create event trigger, the file may not ready, consider using chmod event, and using hashsum the check if the file content is changed.
	task := entity.FileProcessTask{
		File: event.Name,
	}
	bus.Send(&task)
}

func (f *FileSystemWatchTask) Stop(ctx context.Context) error {
	return f.watcher.Close()
}

func (f *FileSystemWatchTask) Append(task server.ITask) {
}
