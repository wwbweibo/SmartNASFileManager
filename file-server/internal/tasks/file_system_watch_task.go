package tasks

import (
	"context"
	domainFile "fileserver/internal/domain/file"
	"fileserver/internal/server"
	"fileserver/utils"
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
	return "file_system_watch_task"
}

func (f *FileSystemWatchTask) GetRunningDuration() time.Duration {
	return 0
}

func (f *FileSystemWatchTask) Start(ctx context.Context) (err error) {
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = f.watcher.Add(f.nasRootPath)
	if err != nil {
		return err
	}
	for {
		select {
		case event := <-f.watcher.Events:
			f.onFileSystemEvent(event)
		case err := <-f.watcher.Errors:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (f *FileSystemWatchTask) onFileSystemEvent(event fsnotify.Event) {
	// do something
	if !f.opts.ShouldWatch(event.Name) {
		return
	}
	switch event.Op {
	case fsnotify.Create:
		f.onFileCreate(event)
	case fsnotify.Remove:
		f.onFileRemove(event)
	}

}

func (f *FileSystemWatchTask) onFileRemove(event fsnotify.Event) {
	if !(utils.CheckIsDir(event.Name)) {
		// remove file from db
		f.repo.RemoveFile(context.Background(), event.Name)
	} else {
		f.repo.RemoveDir(context.Background(), event.Name)
	}
}

func (f *FileSystemWatchTask) onFileCreate(event fsnotify.Event) {
	// 判断是否
}

func (f *FileSystemWatchTask) Stop(ctx context.Context) error {
	return f.watcher.Close()
}

func (f *FileSystemWatchTask) Append(task server.ITask) {
}
