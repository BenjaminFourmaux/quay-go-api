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

func CreateRepositoryTransaction(repository Models.Repository) (*Models.Repository, error) {
	err := Database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the repository
		if err := tx.Create(&repository).Error; err != nil {
			return err // rollback
		}

		// 2. Define Repository Action Count model
		actionCount := Models.RepositoryActionCount{
			RepositoryId: repository.ID,
			Count:        0,
			Date:         time.Now(),
		}

		// 3. Create Repository Action Count entry
		if err := tx.Create(&actionCount).Error; err != nil {
			return err
		}

		// 4. Define Repository Search Score model
		searchScore := Models.RepositorySearchScore{
			RepositoryId: repository.ID,
			Score:        0,
		}

		// 5. Create Repository Search Score entry
		if err := tx.Create(&searchScore).Error; err != nil {
			return err
		}

		return nil // commit
	})

	return &repository, err
}
