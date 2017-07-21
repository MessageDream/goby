package dto

type BlobInfo struct {
	Size uint64 `json:"size"`
	URL  string `json:"url"`
}

type PackageInfo struct {
	AppVersion  string `json:"appVersion"`
	Description string `json:"description"`
	IsDisabled  bool   `json:"isDisabled"`
	IsMandatory bool   `json:"isMandatory"`
	Label       string `json:"label"`
	PackageHash string `json:"packageHash"`
	Rollout     uint8  `json:"rollout"` // unconfirm
}

type Package struct {
	*PackageInfo
	BlobURL            string               `json:"blobUrl"`
	DiffPackageMap     map[string]*BlobInfo `json:"diffPackageMap"`
	OriginalLabel      string               `json:"originalLabel"`
	OriginalDeployment string               `json:"originalDeployment"`
	ReleasedBy         string               `json:"releasedBy"`
	ReleaseMethod      string               `json:"releaseMethod"`
	Size               int64                `json:"size"`
	UploadTime         int64                `json:"uploadTime"`
}

type PackageMetrics struct {
	Active      uint64 `json:"active"`
	TotalActive uint64 `json:"totalActive"`
	Downloaded  uint64 `json:"downloaded"`
	Failed      uint64 `json:"failed"`
	Installed   uint64 `json:"installed"`
}

type UpdatePackageInfo struct {
	PackageInfo *PackageInfo `json:"packageInfo"`
}
