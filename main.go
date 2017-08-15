package main

import (
	"flag"
	"fmt"
	"net/http"
	"octlink/rstore/api"
	"octlink/rstore/configuration"
	"octlink/rstore/handlers"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/utils"
	"octlink/rstore/utils/octlog"
	"os"
)

var (
	port   int
	addr   string
	config string
	cert   string
)

func initDebugConfig() {
	octlog.InitDebugConfig(octlog.DEBUG_LEVEL)
}

func initLogConfig() {
	api.InitApiLog(octlog.DEBUG_LEVEL)
	blobs.InitLog(octlog.DEBUG_LEVEL)
}

func init() {
	flag.StringVar(&config, "config", "var/config.yml", "Config file path")

	initDebugConfig()
	initLogConfig()
}

func usage() {
	fmt.Println("  RVM Store of V" + utils.Version() + "\n")
	fmt.Println("  ./rstore -config ./config.yml\n")
	flag.PrintDefaults()
}

func resolveConfiguration(configfile string) (*configuration.Configuration, error) {

	var configurationPath string

	if configfile == "" {
		configurationPath = os.Getenv("REGISTRY_CONFIGURATION_PATH")
	} else {
		configurationPath = configfile
	}

	if configurationPath == "" {
		return nil, fmt.Errorf("configuration path unspecified")
	}

	fp, err := os.Open(configurationPath)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	config, err := configuration.Parse(fp)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", configurationPath, err)
	}

	return config, nil
}

const (
	HTTP_SERVER = "0.0.0.0"
	HTTP_PORT   = 8001
)

func runApiThread() {

	api := &api.Api{
		Name: "Mirage API Server",
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", HTTP_SERVER, HTTP_PORT),
		Handler:        api.ApiRouter(),
		MaxHeaderBytes: 1 << 20,
	}

	octlog.Warn("Mirage Engine Started\n")

	err := server.ListenAndServe()
	if err != nil {
		octlog.Error("error to listen\n")
	}
}

func main() {

	flag.Usage = usage
	flag.Parse()

	config, err := resolveConfiguration(config)
	if err != nil {
		fmt.Printf("Resolve Configuration Error[%s]\n", err)
		return
	}

	go runApiThread()

	app := handlers.NewApp()

	octlog.Warn("RSTORE Engine Started\n")

	http.ListenAndServe(config.HTTP.Addr, app.Router)
}
