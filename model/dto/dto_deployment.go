package dto

type Deployment struct {
	Key            string          `json:"key"`
	Name           string          `json:"name"`
	Package        *Package        `json:"package"`
	PackageMetrics *PackageMetrics `json:packageMetrics`
}
