package dto

type UpdateCheck struct {
	*PackageInfo
	DownloadURL            string `json:"downloadURL"`
	IsAvailable            bool   `json:"isAvailable"`
	PackageSize            int64  `json:"packageSize"`
	ShouldRunBinaryVersion bool   `json:"shouldRunBinaryVersion"`
	UpdateAppVersion       bool   `json:"updateAppVersion"`
}

type UpdateCheckCache struct {
	originalPackage UpdateCheck
	rollout         int
	rolloutPackage  UpdateCheck
}
