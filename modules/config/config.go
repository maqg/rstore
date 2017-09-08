package config

const (
	// ImageStatusReady for ready state
	ImageStatusReady = "ready"

	// ImageStatusDownloading for downloading state
	ImageStatusDownloading = "downloading"

	//ImageStatusError for error status
	ImageStatusError = "error"
)

const (
	// ImageTypeRootTemplate for root template
	ImageTypeRootTemplate = "RootVolumeTemplate"

	// ImageTypeDataVolume image type of data volume
	ImageTypeDataVolume = "DataVolumeTemplate"

	// ImageTypeIso for iso type
	ImageTypeIso = "ISO"
)

const (
	// OSTypeLinux for os type of linux
	OSTypeLinux = "linux"

	// OSTypeWindows for os type of windows
	OSTypeWindows = "windows"

	// OSTypeMac for os type of mac
	OSTypeMac = "darwin"
)
