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

func GetTagByNameAndRepository(tagName string, repositoryId int) (Models.Tag, error) {
	var tag Models.Tag
	err := Database.DB.
		Preload("Manifest").
		Preload("TagKind").
		Preload("LinkedTag").
		Where("name = ?", tagName).
		Where("repository_id = ?", repositoryId).
		First(&tag).
		Error
	return tag, err
}
