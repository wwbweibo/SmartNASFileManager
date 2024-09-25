package main

import (
	"context"
	"fileserver/internal/biz"
	"fileserver/internal/controllers"
	"fileserver/internal/domain/file"
	"fileserver/internal/server"
	"fileserver/internal/tasks"
	"fileserver/utils"
	"log"

	domainFile "fileserver/internal/domain/file"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancelF := context.WithCancel(context.Background())
	defer cancelF()
	config := &Config{}
	config.Load("config/config.yaml")

	dbconnection := utils.NewDbConnection(config.DBPath)
	fileRepo := file.NewFileRepository(dbconnection)
	taskServer := initTaskServer(*config, fileRepo)
	ginServer := initGinServer(*config, fileRepo)
	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		return taskServer.Start(ctx)
	})
	errGroup.Go(func() error {
		return ginServer.Start(ctx)
	})
	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}

func initTaskServer(config Config,
	fileRepo domainFile.IFileRepository,
) *server.BackendTaskServer {
	taskServer := server.NewBackendTaskServer()
	fileScanOption := tasks.ScanOptions{}
	fileScanOption = fileScanOption.
		OptionRootPath(config.NasRootPath).
		OptionPlainPath(config.ScanOption.Path...).
		OptionRegexPath(config.ScanOption.RegexPath...).
		OptionExtensions(config.ScanOption.Extensions...)
	imageCompressionTask := tasks.NewImageCompressionTask(config.NasRootPath, config.CachePath)
	// register tasks here
	taskServer.RegisterTask(
		imageCompressionTask,
		tasks.NewSysInitBackendTask(fileScanOption, fileRepo, config.DLConfiguration, *imageCompressionTask),
	)
	return taskServer
}

func initGinServer(config Config, fileRepository file.IFileRepository) *server.GinServer {
	server := server.NewGinServer()
	server.UseStatic("/static", config.NasRootPath)
	server.UseStatic("/cache", config.CachePath)
	fileService := biz.NewFilerService(fileRepository)
	fileController := controllers.NewFileApiControllers(fileService)
	server.RegisterController(fileController)
	return server
}
