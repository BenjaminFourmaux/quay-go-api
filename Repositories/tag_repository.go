package Repositories

import (
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetTagsFromRepository(repositoryId int, includeExpired bool) ([]Models.Tag, error) {
	var tags []Models.Tag
	query := Database.DB.
		Preload("Manifest").
		Preload("TagKind").
		Preload("LinkedTag").
		Where("repository_id = ?", repositoryId)

	if !includeExpired {
		query = query.Where("lifetime_end_ms IS NULL")
	}

	err := query.Find(&tags).Error
	return tags, err
}

func GetTagById(tagId int) (*Models.Tag, error) {
	var tag Models.Tag
	err := Database.DB.
		Preload("Manifest").
		Preload("TagKind").
		Preload("LinkedTag").
		Where("id = ?", tagId).
		First(&tag).
		Error
	return &tag, err
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

func UpdateTagManifest(repositoryId int, tagId int, manifestId int) error {
	return Database.DB.
		Model(&Models.Tag{}).
		Where("id = ? AND repository_id = ?", tagId, repositoryId).
		Update("manifest_id", manifestId).
		Error
}

func DeleteTag(tag Models.Tag, nowMs int64) error {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// TODO: (1. Clear pull statistics) ... features not implemented yet

		// TODO: 2. Clean notifications for tag expiry

		// 3. Update lifetime_end_ms with the nowMs
		if err := tx.Model(&tag).Update("lifetime_end_ms", nowMs).Error; err != nil {
			return err
		}
		/*
			Dev notes: Why just update lifetime_end_ms instead of deleting the tag?
			Copilot response: Because we want to keep a record of the tag for historical purposes, and also to prevent re-creation of the same tag name in the future. This is a common practice in systems that require audit trails or historical data retention.
		*/

		return nil // commit
	})
	return err
}
