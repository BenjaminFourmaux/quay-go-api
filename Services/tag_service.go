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

func GetRepositoryTag(repositoryNamespaced string, tagName string, filters map[string]string, currentUser *Auth.AuthenticatedUser) (Dto.Tag, error) {
	logger.Info("[Tag Service] Get Repository Tag")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("Tag: %s", tagName)
	logger.Debug("With filters: %+v", filters)

	// Check if tagName is correct
	if !Common.IsValidTagName(tagName) {
		return Dto.Tag{}, Errors.TagNameInvalid(tagName)
	}

	// Validating filters
	var filterIncludeVulnerabilities bool = false
	if includeVulnerabilities, ok := filters["include_vulnerabilities"]; ok {
		filterIncludeVulnerabilities = includeVulnerabilities == "true"
	}

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.Tag{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.Tag{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Tag{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.Tag{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.Tag{}, err
		}
	}

	// Retrieve the tag
	tagModel, err := Repositories.GetTagByNameAndRepository(tagName, repoExist.ID)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No tag '%s' found in repository '%s'", tagName, repositoryNamespaced)
			return Dto.Tag{}, Errors.TagNotFound(tagName, repositoryNamespaced)
		default:
			logger.Error("Error retrieving tag from database: %s", err.Error())
			return Dto.Tag{}, err
		}
	}

	// Convert model to dto
	tag := Dto.Tag{
		Name:           tagModel.Name,
		Reversion:      tagModel.Reversion,
		StartTs:        time.UnixMilli(tagModel.LifetimeStartMs),
		ManifestDigest: tagModel.Manifest.Digest,
		IsManifestList: false, // TODO: find how to determine if tag is a manifest list
		Size:           *tagModel.Manifest.LayersCompressedSize,
		LastModified:   time.UnixMilli(tagModel.LifetimeStartMs),
	}

	if filterIncludeVulnerabilities {
		vulnerabilities, err := GetVulnerabilityReportForTag(tagModel.ID)
		if err != nil {
			logger.Error("Error retrieving vulnerabilities for tag '%s': %s", tagName, err.Error())
		}

		tag.Vulnerabilities = &vulnerabilities
	}

	return tag, nil
}

func UpdateRepositoryTag(repositoryNamespaced string, tagName string, updateTag Dto.UpdateTag, currentUser *Auth.AuthenticatedUser) (Dto.Tag, error) {
	logger.Info("[Tag Service] Update Repository Tag")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("Tag: %s", tagName)
	logger.Debug("with models: %+v", updateTag)

	// Check if tagName is correct
	if !Common.IsValidTagName(tagName) {
		return Dto.Tag{}, Errors.TagNameInvalid(tagName)
	}

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.Tag{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.Tag{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Tag{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.Tag{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.Tag{}, err
		}
	}

	// Retrieve the tag
	tagModel, err := Repositories.GetTagByNameAndRepository(tagName, repoExist.ID)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No tag '%s' found in repository '%s'", tagName, repositoryNamespaced)
			return Dto.Tag{}, Errors.TagNotFound(tagName, repositoryNamespaced)
		default:
			logger.Error("Error retrieving tag from database: %s", err.Error())
			return Dto.Tag{}, err
		}
	}

	// if features.IMMUTABLE_TAGS... not yet implemented

	// Update Tag Manifest
	if updateTag.ManifestDigest != nil {
		// Check if the given manifest digest exists in this repository
		manifestModel, err := Repositories.GetRepositoryManifestByDigest(repoExist.ID, *updateTag.ManifestDigest)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No manifest '%s' found in repository '%s'", *updateTag.ManifestDigest, repositoryNamespaced)
				return Dto.Tag{}, Errors.ManifestNotFound(*updateTag.ManifestDigest, repositoryNamespaced)
			default:
				logger.Error("Error retrieving manifest from database: %s", err.Error())
				return Dto.Tag{}, err
			}
		}

		// Update the tag associate manifest
		if err = Repositories.UpdateTagManifest(repoExist.ID, tagModel.ID, manifestModel.ID); err != nil {
			logger.Error("Error updating tag manifest in database: %s", err.Error())
			return Dto.Tag{}, err
		} else {
			logger.Info("Tag manifest updated successfully")
		}

	}

	// Update Tag Expiration
	if updateTag.Expiration != nil {
		var expirationDate time.Time
		expirationDate = time.UnixMilli(*updateTag.Expiration)

		if expirationDate.Before(time.Now()) {
			logger.Warning("Expiration date '%s' is in the past", expirationDate)
			return Dto.Tag{}, Errors.InvalidExpirationDate(expirationDate.String())
		}

		// Change repository tag expiration
		if err = Repositories.UpdateRepositoryTagExpiration(repoExist.ID, int(expirationDate.UnixMilli())); err != nil {
			logger.Error("Error updating tag expiration in database: %s", err.Error())
			return Dto.Tag{}, err
		} else {
			logger.Info("Tag expiration updated successfully")
		}
	}

	// Retrieve the updated tag
	updatedTagModel, err := Repositories.GetTagById(tagModel.ID)
	if err != nil {
		logger.Error("Error retrieving updated tag from database: %s", err.Error())
		return Dto.Tag{}, err
	}

	// Convert model to dto
	var updatedTag Dto.Tag
	updatedTag = Dto.Tag{
		Name:           updatedTagModel.Name,
		Reversion:      updatedTagModel.Reversion,
		StartTs:        time.UnixMilli(updatedTagModel.LifetimeStartMs),
		ManifestDigest: updatedTagModel.Manifest.Digest,
		IsManifestList: false, // TODO: find how to determine if tag is a manifest list
		Size:           *updatedTagModel.Manifest.LayersCompressedSize,
		LastModified:   time.UnixMilli(updatedTagModel.LifetimeStartMs),
	}

	return updatedTag, nil
}
