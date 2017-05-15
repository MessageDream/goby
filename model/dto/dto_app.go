package dto

type App struct {
	Collaborators map[string]*Collaborator `json:"collaborators"`
	Deployments   []string                 `json:"deployments"`
	Name          string                   `json:"name"`
	IsOwner       bool                     `json:"-"`
}
