package systemconfig

import (
	"fmt"
	"octlink/rstore/modules/image"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
	"strings"
)

// SystemConfig for system
type SystemConfig struct {
	Version       string `json:"version"`
	RootDirectiry string `json:"rootDirectory"`
	Available     int64  `json:"available"`
	Used          int64  `json:"used"`
	Account       int    `json:"acccount"`
	Capacity      int64  `json:"capacity"`
	Rate          string `json:"rate"`
	Iso           int    `json:"iso"`
	File          int    `json:"file"`
	Snapshot      int    `json:"snapshot"`
	Root          int    `json:"root"`
	DataDisk      int    `json:"dataDisk"`
}

// GetCapacity for Capacity fetching
func (sc *SystemConfig) GetCapacity() {

	if utils.IsPlatformWindows() {
		return
	}

	cmd := fmt.Sprintf("df %s --total | grep total | grep -v grep | awk -F' ' '{print$2,$3,$4,$5}'",
		configuration.GetConfig().RootDirectory)
	data, err := utils.OCTSystem(cmd)
	if err != nil {
		octlog.Error("exec cmd [%s] error [%s]\n", cmd, err)
		return
	}

	segs := strings.Split(strings.Replace(data, "\n", "", -1), " ")
	sc.Capacity = utils.StringToInt64(segs[0])
	sc.Used = utils.StringToInt64(segs[1])
	sc.Available = utils.StringToInt64(segs[2])
	sc.Rate = segs[3]
}

// GetSystemConfig get system config of this backupstorage
func GetSystemConfig() *SystemConfig {

	conf := configuration.GetConfig()

	sc := new(SystemConfig)
	sc.Version = conf.Version
	sc.RootDirectiry = conf.RootDirectory
	sc.Iso = len(image.GImagesIsoMap)
	sc.Root = len(image.GImagesRootTemplateMap)
	sc.DataDisk = len(image.GImagesDataTemplateMap)

	sc.GetCapacity()

	return sc
}
