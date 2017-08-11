package main

import (
	"flag"
	"fmt"
	"net/http"
	"octlink/mirage/src/api"
	"octlink/mirage/src/utils"
	"octlink/mirage/src/utils/octlog"
	"octlink/rstore/handlers"
	"octlink/rstore/modules/blobs"
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
	flag.StringVar(&config, "config", "./config.json", "config file")

	initDebugConfig()
	initLogConfig()
}

func usage() {
	fmt.Println("  RVM Store of V" + utils.Version() + "\n")
	fmt.Println("  ./rstore -config ./config.json\n")
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	app := handlers.NewApp()

	http.ListenAndServe(":8000", app.Router)

	octlog.Warn("RSTORE Engine Started\n")
}
