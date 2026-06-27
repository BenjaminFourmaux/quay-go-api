package Dto

import "time"

type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"` // Image name
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Kind        string `json:"kind"` // Image or Application
	State       string `json:"state"`
	IsStarred   bool   `json:"is_starred"`
}

type CreateRepository struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Namespace   *string `json:"namespace"` // Can be null if is a global repo, not organization scoped
	IsPublic    bool    `json:"is_public"`
	Kind        string  `json:"kind"`
}

type RepositoryDetails struct {
	Namespace      string            `json:"namespace"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Kind           string            `json:"kind"`
	IsPublic       bool              `json:"is_public"`
	IsOrganization bool              `json:"is_organization"`
	IsStarred      bool              `json:"is_starred"`
	StatusToken    string            `json:"status_token"`
	TrustEnabled   bool              `json:"trust_enabled"`
	TagExpirationS int               `json:"tag_expiration_s"`
	State          string            `json:"state"`
	Tags           []RepositoryTag   `json:"tags,omitempty"`
	CanWrite       bool              `json:"can_write"`
	CanAdmin       bool              `json:"can_admin"`
	Stats          []RepositoryStats `json:"stats,omitempty"`
}

type RepositoryTag struct {
	Name           string    `json:"name"`
	Size           int64     `json:"size"` // In bits
	LastModified   time.Time `json:"last_modified"`
	ManifestDigest string    `json:"manifest_digest"`
}

type RepositoryStats struct {
	Date  time.Time `json:"date"` // In format: yyyy-mm-dd
	Count int       `json:"count"`
}
