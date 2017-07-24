package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/mirage/src/utils/octlog"
	"os"
	"strings"
)

const (
	PARAM_TYPE_STRING     = "string"
	PARAM_TYPE_INT        = "int"
	PARAM_TYPE_LISTINT    = "listint"
	PARAM_TYPE_LISTSTRING = "liststring"
	PARAM_TYPE_BOOLEAN    = "boolean"

	PARAM_NOT_NULL = "NotNull"

	API_PREFIX_CENTER = "octlink.mirage.center"
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
	Timeout int         `json:"timeout"`
	Name    string      `json:"name"`
	Key     string      `json:"key"`
	Handler string      `json:"handler"`
	Paras   []ProtoPara `json:"paras"`
}

func InitApiLog(level int) {
	logger = octlog.InitLogConfig("api.log", level)
}

type ApiModule struct {
	Name   string              `json:"name"`
	Protos map[string]ApiProto `json:"protos"`
}

// octlink.mirage.center.host.APIAddHost
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

func loadModule(name string, baseDir string) bool {

	var apiModule ApiModule

	filePath := baseDir + "apiconfig/" + name + ".json"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file " + filePath + "error")
		return false
	}

	data, err := ioutil.ReadFile(filePath)

	file.Close()

	err = json.Unmarshal(data, &apiModule.Protos)
	if err != nil {
		fmt.Println("Transfer json bytes error")
		fmt.Println(err)
		return false
	}
	apiModule.Name = name

	if GApiConfig.Modules == nil {
		GApiConfig.Modules = make(map[string]ApiModule, 50)
	}

	for key, proto := range apiModule.Protos {
		proto.Key = API_PREFIX_CENTER + "." + name + "." + key
		apiModule.Protos[key] = proto
	}

	GApiConfig.Modules[name] = apiModule

	return true
}

func LoadApiConfig(baseDir string) bool {

	modules := []string{"user", "host", "account", "usergroup"}

	for i := 0; i < len(modules); i++ {
		state := loadModule(modules[i], baseDir)
		if state != true {
			fmt.Println("load module %s error", modules[i])
			return false
		}
	}

	InitApiService()

	return true
}
