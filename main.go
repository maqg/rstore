package main

import (
	"flag"
	"fmt"
	"net/http"
	"octlink/mirage/src/api"
	"octlink/mirage/src/utils"
	"octlink/mirage/src/utils/octlog"
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
}

func init() {
	flag.StringVar(&config, "config", "./config.json", "config file")
}

func usage() {
	fmt.Println("  RVM Store of V" + utils.Version() + "\n")
	fmt.Println("  ./rstore -config ./config.json\n")
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	fmt.Println(utils.Version())

	initDebugConfig()
	initLogConfig()

	api := &api.Api{
		Name: "RSTORE API Server",
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", addr, port),
		Handler:        api.ApiRouter(),
		MaxHeaderBytes: 1 << 20,
	}

	octlog.Warn("RSTORE Engine Started\n")

	err := server.ListenAndServe()
	if err != nil {
		octlog.Error("error to listen\n")
	}
}
