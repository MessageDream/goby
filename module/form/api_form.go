package form

type AppOrDeploymentOption struct {
	Name string `json:"name" binding:"Required"`
}

type PackageInfo struct {
	AppVersion  *string `json:"appVersion"`
	Description *string `json:"description"`
	IsDisabled  *bool   `json:"isDisabled"`
	IsMandatory *bool   `json:"isMandatory"`
	Label       *string `json:"label"`
	Rollout     *uint8  `json:"rollout"`
}

type UpdatePackageInfo struct {
	PackageInfo *PackageInfo `json:"packageInfo"`
}

type ReportStatus struct {
	AppVersion                *string `json:"appVersion"`
	DeploymentKey             *string `json:"deploymentKey"`
	Label                     *string `json:"label"`
	Status                    *string `json:"status"`
	ClientUniqueId            *string `json:"clientUniqueId"`
	PreviousDeploymentKey     *string `json:"previousDeploymentKey"`
	PreviousLabelOrAppVersion *string `json:"previousLabelOrAppVersion"`
}
