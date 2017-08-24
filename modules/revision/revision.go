package revision

type Revision struct {
	Id          string `json:"uuid"`
	Name        string `json:"name"`
	BlobSum     string `json:"blobsum"`
	CraeteTime  string `json:"createTime"`
	DiskSize    int64  `json:"diskSize"`
	VirtualSize int64  `json:"virtualSize"`
}
