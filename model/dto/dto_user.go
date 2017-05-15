package dto

import "time"

type UserToken struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	FriendlyName string    `json:"friendlyName"`
	Expires      time.Time `json:"expires"`
	IsSession    bool      `json:"isSession"`
	CreatedBy    string    `json:"createdBy"`
	CreatedTime  time.Time `json:"createdTime"`
	Description  string    `json:"description"`
}

type User struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	LinkedProviders []string `json:"linkedProviders"`
}

type Collaborator struct {
	IsCurrentAccount bool   `json:"isCurrentAccount"`
	Permission       string `json:"permission"`
}
