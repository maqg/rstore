package main

import "octlink/rstore/cmd/cli/clis"

/*
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
*/

func main() {
	clis.RootCmd.Execute()
}
