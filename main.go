package main

import (
	"flag"
	"fmt"
	"net/http"
	"octlink/rstore/api"
	"octlink/rstore/handlers"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/image"
	"octlink/rstore/modules/manifest"
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

func initDebugConfig() {
	octlog.InitDebugConfig(octlog.DebugLevel)
}

func initLogConfig() {
	api.InitAPILog(octlog.DebugLevel)
	blobs.InitLog(octlog.DebugLevel)
}

func init() {
	flag.StringVar(&config, "config", "./config.yml", "Config file path")

	initDebugConfig()
	initLogConfig()
}

func usage() {
	fmt.Printf("  RVM Store of V" + utils.Version() + "\n")
	fmt.Printf("  ./rstore -config ./config.yml\n")
	flag.PrintDefaults()
}

func runAPIThread(conf *configuration.Configuration) {

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
		octlog.Error("error to listen\n")
	}
}

func initRootDirectory(conf *configuration.Configuration) {
	utils.CreateDir(conf.RootDirectory + manifest.ReposDir)
	utils.CreateDir(conf.RootDirectory + manifest.BlobDir)
	utils.CreateDir(conf.RootDirectory + manifest.BlobManifestDir)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	conf, err := configuration.ResolveConfig(config)
	if err != nil {
		fmt.Printf("Resolve Configuration Error[%s]\n", err)
		return
	}

	initRootDirectory(conf)

	// ReloadImages here
	image.ReloadImages()

	go runAPIThread(conf)

	app := handlers.NewApp()

	octlog.Warn("RSTORE HTTP Engine Started ON %s\n", conf.HTTP.Addr)

	http.ListenAndServe(conf.HTTP.Addr, app.Router)
}
