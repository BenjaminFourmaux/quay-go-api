package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
	"time"
)

func ListRepositoryTags(repositoryNamespaced string, filters map[string]string, currentUser *Auth.AuthenticatedUser) ([]Dto.Tag, error) {
	logger.Info("[Tag Service] List Repository Tags")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var filterIncludeVulnerabilities bool = false
	if includeVulnerabilities, ok := filters["include_vulnerabilities"]; ok {
		filterIncludeVulnerabilities = includeVulnerabilities == "true"
	}

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return []Dto.Tag{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return []Dto.Tag{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.Tag{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return []Dto.Tag{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return []Dto.Tag{}, err
		}
	}

	// Retrieve the list of tags in this repository
	tagsModel, err := Repositories.GetTagsFromRepository(repoExist.ID)
	if err != nil {
		logger.Error("Error retrieving tags from repository: %s", err.Error())
		return []Dto.Tag{}, err
	}

	// Convert models to dto
	var tags []Dto.Tag
	for _, tagModel := range tagsModel {
		tag := Dto.Tag{
			Name:           tagModel.Name,
			Reversion:      tagModel.Reversion,
			StartTs:        time.UnixMilli(tagModel.LifetimeStartMs),
			ManifestDigest: tagModel.Manifest.Digest,
			IsManifestList: false, // TODO: find how to determine if tag is a manifest list
			Size:           *tagModel.Manifest.LayersCompressedSize,
			LastModified:   time.UnixMilli(tagModel.LifetimeStartMs),
		}

		// Include vulnerability information if requested
		if filterIncludeVulnerabilities {
			vulnerabilities, err := GetVulnerabilityReportForTag(tagModel.ID)
			if err != nil {
				logger.Error("Error retrieving vulnerabilities for tag '%s': %s", tagModel.Name, err.Error())
				return []Dto.Tag{}, err
			}
			tag.Vulnerabilities = &vulnerabilities
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
