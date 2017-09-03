package config

import "octlink/rstore/utils/configuration"

const (
	// ImageStatusReady for ready state
	ImageStatusReady = "ready"

	// ImageStatusDownloading for downloading state
	ImageStatusDownloading = "downloading"

	//ImageStatusError for error status
	ImageStatusError = "error"
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

// GetSystemConfig get system config of this backupstorage
func GetSystemConfig() *SystemConfig {

	conf := configuration.GetConfig()

	sc := new(SystemConfig)
	sc.Version = conf.Version
	sc.RootDirectiry = conf.RootDirectory

	return sc
}
