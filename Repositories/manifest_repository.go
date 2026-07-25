package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetRepositoryManifestByDigest(repositoryId int, digest string) (Models.Manifest, error) {
	var manifest Models.Manifest
	err := Database.DB.
		Where("repository_id = ?", repositoryId).
		Where("digest = ?", digest).
		First(&manifest).Error
	return manifest, err
}
