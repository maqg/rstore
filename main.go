package main

import (
	"flag"
	"fmt"
	"net/http"
	"octlink/rstore/api"
	"octlink/rstore/handlers"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/image"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/modules/task"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
)

var (
	port   int
	addr   string
	config string
	cert   string
)

var conf *configuration.Configuration

func initDebugConfig() {
	octlog.InitDebugConfig(conf.DebugLevel)
}

func initLogConfig() {

	utils.CreateDir(conf.RootDirectory + conf.LogDirectory)

	api.InitAPILog(conf.LogLevel)

	blobs.InitLog(conf.LogLevel)

	blobsmanifest.InitLog(conf.LogLevel)

	utils.InitLog(conf.LogLevel)

	image.InitLog(conf.LogLevel)

	handlers.InitLog(conf.LogLevel)

	task.InitLog(conf.LogLevel)
}

func initDebugAndLog() {
	initDebugConfig()
	initLogConfig()
}

func init() {
	flag.StringVar(&config, "config", "./config.yml", "Config file path")
}

func usage() {
	fmt.Printf("  RVM Store of V" + utils.Version() + "\n")
	fmt.Printf("  ./rstore -config ./config.yml\n")
	flag.PrintDefaults()
}

func runAPIThread() {

	api := &api.API{
		Name: "Rstore API Server",
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s", conf.HTTP.APIAddr),
		Handler:        api.Router(),
		MaxHeaderBytes: 1 << 20,
	}

	octlog.Warn("RSTORE API Engine Started ON %s\n", conf.HTTP.APIAddr)

	err := server.ListenAndServe()
	if err != nil {
		octlog.Error("error to listen at %s\n", conf.HTTP.APIAddr)
	}
}

func initRootDirectory() {
	utils.CreateDir(conf.RootDirectory + manifest.ReposDir)
	utils.CreateDir(conf.RootDirectory + manifest.BlobDir)
	utils.CreateDir(conf.RootDirectory + manifest.BlobManifestDir)
	utils.CreateDir(conf.RootDirectory + manifest.TempDir)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	c, err := configuration.ResolveConfig(config)
	if err != nil {
		fmt.Printf("Resolve Configuration Error[%s]\n", err)
		return
	}
	conf = c

	// for root direcotry
	initRootDirectory()

	// for debug and log config
	initDebugAndLog()

	// ReloadImages here
	image.ReloadImages()
	go image.ReloadSignal()

	go runAPIThread()

	app := handlers.NewApp()

	octlog.Warn("RSTORE HTTP Engine Started ON %s\n", conf.HTTP.Addr)

	err = http.ListenAndServe(conf.HTTP.Addr, app.Router)
	if err != nil {
		fmt.Printf("error to listen at data address %s\n", conf.HTTP.Addr)
	}
}
