package main

import (
	"octlink/rstore/api"
	"octlink/rstore/cmd/cli/clis"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/utils/octlog"
)

func initDebugConfig() {
	octlog.InitDebugConfig(octlog.DebugLevel)
}

func initLogConfig() {
	api.InitAPILog(octlog.DebugLevel)
	blobs.InitLog(octlog.DebugLevel)
}

func init() {
	initDebugConfig()
	initLogConfig()
}

func main() {
	clis.RootCmd.Execute()
}
