package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
)

/*
GetManifest Not fully supported function, need to call official API to getting information
*/
func GetManifest(repositoryNamespaced string, manifestRef string, currentUser *Auth.AuthenticatedUser) (Dto.Manifest, error) {
	logger.Info("[Manifest Service] Get Manifest")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("Manifest ref: %s", manifestRef)

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.Manifest{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.Manifest{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Manifest{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.Manifest{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.Manifest{}, err
		}
	}

	// Get the manifest and check if exists
	manifestModel, err := Repositories.GetRepositoryManifestByDigest(repoExist.ID, manifestRef)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No manifest '%s' found in repository '%s'", manifestRef, repositoryNamespaced)
			return Dto.Manifest{}, Errors.ManifestNotFound(manifestRef, repositoryNamespaced)
		default:
			logger.Error("Error retrieving manifest from database: %s", err.Error())
			return Dto.Manifest{}, err
		}
	}

	// Get layers informations
	// TODO: we need to call the official API to get layers information, because the current implementation does not support it (access to blob storage).
	quayManifestModel, err := QuayGetManifest(repositoryNamespaced, manifestRef, *currentUser)
	if err != nil {
		logger.Warning("Error retrieving manifest from Quay API: %s, pass...", err.Error())
		// pass
		quayManifestModel = Dto.Manifest{
			Layers: []Dto.Layer{},
		}
	}

	// Convert manifest model to DTO
	manifestDTO := Dto.Manifest{
		Digest:          manifestModel.Digest,
		IsManifestList:  false, // TODO: find how to determine if tag is a manifest list
		ManifestData:    manifestModel.ManifestBytes,
		ConfigMediaType: manifestModel.ConfigMediaType,
		Layers:          quayManifestModel.Layers,
	}

	return manifestDTO, nil
}
