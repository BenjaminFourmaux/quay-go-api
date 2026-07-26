package Dto

import "time"

type Manifest struct {
	Digest          string          `json:"digest"`
	IsManifestList  bool            `json:"is_manifest_list"`
	ManifestData    string          `json:"manifest_data"`
	ConfigMediaType *string         `json:"config_media_type"`
	Layers          []ManifestLayer `json:"layers"`
}

type ManifestLayer struct {
	Index           int       `json:"index"`
	CompressedSize  int       `json:"compressed_size"` // In Bytes
	IsRemote        bool      `json:"is_remote"`
	Urls            *[]string `json:"urls"`
	Command         []string  `json:"command"`
	Comment         string    `json:"comment"`
	Author          *string   `json:"author"`
	BlobDigest      string    `json:"blob_digest"`
	CreatedDatetime time.Time `json:"created_datetime"`
}

type ManifestLabel struct {
	Id         string `json:"id"` // UUID
	Key        string `json:"key"`
	Value      string `json:"value"`
	SourceType string `json:"source_type"`
	MediaType  string `json:"media_type"`
}
