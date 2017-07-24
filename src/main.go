package main

import (
	"fmt"
	"net/http"
	"octlink/mirage/src/api"
	"octlink/mirage/src/modules/account"
	"octlink/mirage/src/modules/session"
	"octlink/mirage/src/modules/user"
	"octlink/mirage/src/modules/usergroup"
	"octlink/mirage/src/utils"
	"octlink/mirage/src/utils/octlog"
)

func initDebugConfig() {
	octlog.InitDebugConfig(octlog.DEBUG_LEVEL)
}

func initLogConfig() {

	api.InitApiLog(octlog.DEBUG_LEVEL)

	user.InitLog(octlog.DEBUG_LEVEL)
	account.InitLog(octlog.DEBUG_LEVEL)
	usergroup.InitLog(octlog.DEBUG_LEVEL)
	session.InitLog(octlog.DEBUG_LEVEL)
}

const (
	HTTP_SERVER = "0.0.0.0"
	HTTP_PORT   = 8080
)

func main() {

	fmt.Println(utils.Version())

	initDebugConfig()
	initLogConfig()

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
