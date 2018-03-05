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
	// ZeroDataDigest8M 8 Mb bytes zero data digest value
	ZeroDataDigest8M = "2daeb1f36095b44b318410b3f4e8b5d989dcc7bb023d1426c492dab0a3053e74"

	// ZeroDataDigest4M for 4 Mb zero bytes digest value
	ZeroDataDigest4M = "bb9f8df61474d25e71fa00722318cd387396ca1736605e1248821cc0de3d3af8"
)

// ZeroData8M for zero data 8 MB
var ZeroData8M = make([]byte, 8*1024*1024)

// ZeroData4M for zero data of 4MB
var ZeroData4M = make([]byte, 4*1024*1024)

const (
	// OSTypeLinux for os type of linux
	OSTypeLinux = "linux"

	// OSTypeWindows for os type of windows
	OSTypeWindows = "windows"

	// OSTypeMac for os type of mac
	OSTypeMac = "darwin"
)
