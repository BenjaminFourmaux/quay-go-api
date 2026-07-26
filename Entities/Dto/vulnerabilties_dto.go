package Dto

import "time"

type Vulnerabilities struct {
	Critical    int `json:"critical"`
	High        int `json:"high"`
	Medium      int `json:"medium"`
	Low         int `json:"low"`
	Unspecified int `json:"unspecified"`
}

type VulnerabilityCVE struct {
	CVEID            string    `json:"cve_id"`
	SeverityScore    int       `json:"severity_score"`
	Fixable          string    `json:"fixable"`    // version in witch this vulnerability is fixed, if empty then it is not fixable
	PresentIn        string    `json:"present_in"` // 'layer' or 'base'
	AffectedPackages string    `json:"affected_packages"`
	AffectedVersions string    `json:"affected_versions"`
	PublishedDate    time.Time `json:"published_date"`
}
