package Dto

import (
	"fmt"
	"time"
)

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

type UpdateTag struct {
	Immutable      bool    `json:"immutable"`       // (If specified) Whether the tag should be immutable. Write permission required to set, admin permission required to unset.
	Expiration     *int64  `json:"expiration"`      // (If specified) The expiration for the image
	ManifestDigest *string `json:"manifest_digest"` // (If specified) The manifest digest to which the tag should point
}

func (u UpdateTag) String() string {
	exp := "<nil>"
	if u.Expiration != nil {
		exp = fmt.Sprintf("%d", *u.Expiration)
	}

	digest := "<nil>"
	if u.ManifestDigest != nil {
		digest = *u.ManifestDigest
	}

	return fmt.Sprintf("{Immutable:%t Expiration:%s ManifestDigest:%s}", u.Immutable, exp, digest)
}
