package Dto

import "time"

type Tag struct {
	Name           string    `json:"name"`
	Reversion      bool      `json:"reversion"`
	StartTs        time.Time `json:"start_ts"`
	ManifestDigest string    `json:"manifest_digest"`
	IsManifestList bool      `json:"is_manifest_list"`
	Size           int64     `json:"size"` // In bits
	LastModified   time.Time `json:"last_modified"`

	// Additional fields for vulnerability information
	Vulnerabilities *Vulnerabilities `json:"vulnerabilities,omitempty"`
}
