package main

import (
	"context"
	"fileserver/file"
	"fileserver/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx, cancelF := context.WithCancel(context.Background())
	defer cancelF()
	config := &Config{}
	config.Load("config/config.yaml")
	fileScanOption := file.ScanOptions{}
	fileScanOption = fileScanOption.OptionPlainPath(config.ScanOption.Path...).
		OptionRegexPath(config.ScanOption.RegexPath...).
		OptionExtensions(config.ScanOption.Extensions...)

	dbconnection := utils.NewDbConnection()
	fileRepo := file.NewFileRepository(dbconnection)
	fileChan := make(chan string, 10)
	go file.StartFileScanner(ctx, fileChan, fileRepo)
	go file.ScanAndUpdateFiles(ctx, config.NasRootPath, fileScanOption, fileChan)
	// start a http server here
	http.HandleFunc("/", handler)
	// handle for static files
	http.Handle("/nas/", http.FileServer(http.Dir(config.NasRootPath)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
