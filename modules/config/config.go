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
