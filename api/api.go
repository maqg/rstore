package api

import (
	"fmt"
	"octlink/rstore/utils/octlog"
	"strings"
)

const (
	PARAM_TYPE_STRING     = "string"
	PARAM_TYPE_INT        = "int"
	PARAM_TYPE_LISTINT    = "listint"
	PARAM_TYPE_LISTSTRING = "liststring"
	PARAM_TYPE_BOOLEAN    = "boolean"

	PARAM_NOT_NULL = "NotNull"

	API_PREFIX_CENTER = "octlink.rstore.center"
)

var logger *octlog.LogConfig

var GApiConfig ApiConfig

type Api struct {
	Name   string
	Config *ApiConfig
}

type ApiConfig struct {
	Modules map[string]ApiModule `json:"modules`
}

type ProtoPara struct {
	Name    string      `json:"name"`
	Default interface{} `json:"default"`
	Type    string      `json:"type"`
	Desc    string      `json:"desc"`
}

type ApiProto struct {
	Name    string      `json:"name"`
	Key     string      `json:"key"`
	Paras   []ProtoPara `json:"paras"`
	handler func(*ApiParas) *ApiResponse
}

func InitApiLog(level int) {
	logger = octlog.InitLogConfig("api.log", level)
}

type ApiModule struct {
	Name   string              `json:"name"`
	Protos map[string]ApiProto `json:"protos"`
}

// octlink.rstore.center.host.APIAddHost
func FindApiProto(api string) *ApiProto {

	segments := strings.Split(api, ".")
	moduleName := segments[3]
	apiKey := segments[4]

	if moduleName == "" || apiKey == "" {
		fmt.Printf("got bad api key %s\n", api)
		return nil
	}

	module, ok := GApiConfig.Modules[moduleName]
	if !ok {
		fmt.Printf("no module exist for %s\n", moduleName)
		return nil
	}

	proto, ok := module.Protos[apiKey]
	if !ok {
		fmt.Printf("no proto exist for %s\n", apiKey)
		return nil
	}

	return &proto
}
