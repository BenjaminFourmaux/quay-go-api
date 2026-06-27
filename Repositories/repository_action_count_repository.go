package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetCountsFromRepository(repositoryId int) ([]Models.RepositoryActionCount, error) {
	counts := make([]Models.RepositoryActionCount, 0)
	err := Database.DB.Find(&counts, "repository_id = ?", repositoryId).Error
	return counts, err
}
