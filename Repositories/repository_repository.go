package Repositories

import (
	"gorm.io/gorm"
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
	"time"
)

func SelectRepositories(userId int, organizationOrUserId *int) ([]Models.Repository, error) {
	var repositories []Models.Repository
	query := Database.DB.
		Preload("NamespaceUser").
		Preload("Visibility").
		Preload("Kind").
		Preload("Stars", "user_id = ?", userId)

	if organizationOrUserId != nil {
		query = query.Where("namespace_user_id = ?", *organizationOrUserId)
	}

	err := query.Find(&repositories).Error
	return repositories, err
}

func GetOrganizationRepositoriesByOrgId(organizationId int, userId int) ([]Models.Repository, error) {
	var repositories []Models.Repository
	err := Database.DB.
		Preload("Visibility").
		Preload("Kind").
		Preload("Stars", "user_id = ?", userId).
		Where("namespace_user_id = ?", organizationId).
		Find(&repositories).
		Error
	return repositories, err
}

func GetRepositoryById(repositoryId int, userId int) (Models.Repository, error) {
	var repository Models.Repository
	err := Database.DB.
		Preload("NamespaceUser").
		Preload("Visibility").
		Preload("Kind").
		Preload("Stars", "user_id = ?", userId).
		First(&repository, repositoryId).
		Error
	return repository, err
}

func FindRepositoryByNameAndNamespace(name string, namespace *string) (Models.Repository, error) {
	var repository Models.Repository
	query := Database.DB.
		Model(&Models.Repository{}).
		Where("repository.name = ?", name)

	if namespace != nil {
		query = query.
			Joins("NamespaceUser", Database.DB.Where(&Models.User{Username: *namespace}))
	} else {
		query = query.Where("namespace_user_id IS NULL")
	}

	err := query.First(&repository).Error
	return repository, err
}

func CreateRepositoryTransaction(repository Models.Repository, userId int) (*Models.Repository, error) {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the repository
		if err := tx.Create(&repository).Error; err != nil {
			return err // rollback
		}

		// 2. Define Repository Permission to add the user as admin
		permission := Models.RepositoryPermission{
			UserId:       &userId,
			RepositoryId: repository.ID,
			RoleId:       1, // Admin role
		}

		// 3. Create the Repository Permission entry
		if err := tx.Create(&permission).Error; err != nil {
			return err // rollback
		}

		// 4. Define Repository Action Count model
		actionCount := Models.RepositoryActionCount{
			RepositoryId: repository.ID,
			Count:        0,
			Date:         time.Now(),
		}

		// 5. Create Repository Action Count entry
		if err := tx.Create(&actionCount).Error; err != nil {
			return err
		}

		// 6. Define Repository Search Score model
		searchScore := Models.RepositorySearchScore{
			RepositoryId: repository.ID,
			Score:        0,
		}

		// 7. Create Repository Search Score entry
		if err := tx.Create(&searchScore).Error; err != nil {
			return err
		}

		return nil // commit
	})

	return &repository, err
}

func DeleteRepositoryTransaction(repository Models.Repository) error {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Remove tags linked to the repository
		if err := tx.Where("repository_id = ?", repository.ID).Delete(&Models.Tag{}).Error; err != nil {
			return err
		}

		// 2. Remove manifests linked to the repository
		if err := tx.Where("repository_id = ?", repository.ID).Delete(&Models.Manifest{}).Error; err != nil {
			return err
		}

		// 3. Remove stars linked to the repository
		if err := tx.Where("repository_id = ?", repository.ID).Delete(&Models.Star{}).Error; err != nil {
			return err
		}

		// 4. Remove permissions linked to the repository
		if err := tx.Where("repository_id = ?", repository.ID).Delete(&Models.RepositoryPermission{}).Error; err != nil {
			return err
		}

		// 5. Remove action counts linked to the repository
		if err := tx.Where("repository_id = ?", repository.ID).Delete(&Models.RepositoryActionCount{}).Error; err != nil {
			return err
		}

		// 6. Remove search score linked to the repository
		if err := tx.Where("repository_id = ?", repository.ID).Delete(&Models.RepositorySearchScore{}).Error; err != nil {
			return err
		}

		// 7. Remove the repository
		deleteResult := tx.Where("id = ?", repository.ID).Delete(&Models.Repository{})
		if deleteResult.Error != nil {
			return deleteResult.Error
		}

		if deleteResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	return err
}

func UpdateRepository(repository Models.Repository) (*Models.Repository, error) {
	err := Database.DB.Preload("Kind").Preload("NamespaceUser").Save(&repository).Error
	return &repository, err
}
