package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetTagsFromRepository(repositoryId int) ([]Models.Tag, error) {
	var tags []Models.Tag
	err := Database.DB.
		Preload("Manifest").
		Preload("TagKind").
		Preload("LinkedTag").
		Where("repository_id = ?", repositoryId).
		Find(&tags).
		Error
	return tags, err
}
