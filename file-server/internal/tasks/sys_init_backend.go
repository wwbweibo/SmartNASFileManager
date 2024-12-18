package tasks

import (
	"context"
	"fileserver/internal/adapters/dl"
	domainFile "fileserver/internal/domain/file"
	"fileserver/internal/server"
	"fileserver/internal/tasks/entity"
	"fileserver/utils"
	"log"
	"strings"
	"time"
)

type SysInitBackendTask struct {
	startTime time.Time
	option    utils.ScanOptions
	repo      domainFile.IFileRepository
	dlConfig  dl.Config
}

func NewSysInitBackendTask(option utils.ScanOptions,
	repo domainFile.IFileRepository,
	dlConfig dl.Config,
) *SysInitBackendTask {
	return &SysInitBackendTask{
		option:   option,
		repo:     repo,
		dlConfig: dlConfig,
	}
}

func (s *SysInitBackendTask) GetTaskName() string {
	return "SysInitBackendTask"
}

func (s *SysInitBackendTask) GetRunningDuration() time.Duration {
	return time.Duration(0)
}

func (s *SysInitBackendTask) Start(ctx context.Context) error {
	s.startTime = time.Now()
	files, _ := utils.WalkDir(s.option.RootPath)
	log.Default().Printf("found %d files", len(files))
	for _, file := range files {
		select {
		case <-ctx.Done():
			return nil
		default:
			if s.option.ShouldWatch(file, false) {
				domainFile.Root.Add(utils.GetDirectory(strings.Replace(file, s.option.RootPath, "", 1)))
				bus.Send(&entity.FileProcessTask{
					File: file,
				})
			}
		}
	}
	return nil
}

func (s *SysInitBackendTask) Stop(ctx context.Context) error {
	return nil
}

func (s *SysInitBackendTask) Append(task server.ITask) {
}
