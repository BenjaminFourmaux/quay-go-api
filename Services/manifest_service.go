package Services

import (
	"github.com/google/uuid"
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
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
			Layers: []Dto.ManifestLayer{},
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

func GetManifestLabels(repositoryNamespaced string, manifestRef string, currentUser *Auth.AuthenticatedUser) ([]Dto.ManifestLabel, error) {
	logger.Info("[Manifest Service] Get Manifest Labels")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("Manifest ref: %s", manifestRef)

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return []Dto.ManifestLabel{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return []Dto.ManifestLabel{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.ManifestLabel{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return []Dto.ManifestLabel{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return []Dto.ManifestLabel{}, err
		}
	}

	// Get the manifest and check if exists
	manifestModel, err := Repositories.GetRepositoryManifestByDigest(repoExist.ID, manifestRef)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No manifest '%s' found in repository '%s'", manifestRef, repositoryNamespaced)
			return []Dto.ManifestLabel{}, Errors.ManifestNotFound(manifestRef, repositoryNamespaced)
		default:
			logger.Error("Error retrieving manifest from database: %s", err.Error())
			return []Dto.ManifestLabel{}, err
		}
	}

	// Get labels from the database
	manifestLabelsModel, err := Repositories.ListManifestLabels(repoExist.ID, manifestModel.ID)
	if err != nil {
		logger.Error("Error retrieving manifest labels from database: %s", err.Error())
		return []Dto.ManifestLabel{}, err
	}

	// Convert model to DTO
	labels := []Dto.ManifestLabel{}
	for _, labelModel := range manifestLabelsModel {
		labels = append(labels, Dto.ManifestLabel{
			Id:         labelModel.Label.UUID,
			Key:        labelModel.Label.Key,
			Value:      labelModel.Label.Value,
			SourceType: Common.MapLabelSourceType(labelModel.Label.SourceTypeId).Name,
			MediaType:  Common.MapMediaTypeName(labelModel.Label.MediaTypeId),
		})
	}

	return labels, nil
}

func CreateManifestLabel(repositoryNamespaced string, manifestRef string, addLabel Dto.AddManifestLabel, currentUser *Auth.AuthenticatedUser) (Dto.ManifestLabel, error) {
	logger.Info("[Manifest Service] Create Manifest Label")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("Manifest ref: %s", manifestRef)
	logger.Debug("Label %s=%s", addLabel.Key, addLabel.Value)

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.ManifestLabel{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.ManifestLabel{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.ManifestLabel{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.ManifestLabel{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.ManifestLabel{}, err
		}
	}

	// Get the manifest and check if exists
	manifestModel, err := Repositories.GetRepositoryManifestByDigest(repoExist.ID, manifestRef)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No manifest '%s' found in repository '%s'", manifestRef, repositoryNamespaced)
			return Dto.ManifestLabel{}, Errors.ManifestNotFound(manifestRef, repositoryNamespaced)
		default:
			logger.Error("Error retrieving manifest from database: %s", err.Error())
			return Dto.ManifestLabel{}, err
		}
	}

	// Create model to insert
	createLabelModel := Models.Label{
		UUID:         uuid.New().String(),
		Key:          addLabel.Key,
		Value:        addLabel.Value,
		MediaTypeId:  1, // 1 => text/plain
		SourceTypeId: 2, // 2 => api
	}

	// Insert into the database
	createdLabelModel, err := Repositories.AddManifestLabel(repoExist.ID, manifestModel.ID, createLabelModel)
	if err != nil {
		logger.Error("Error inserting manifest label into database: %s", err.Error())
		return Dto.ManifestLabel{}, err
	}

	// Convert model into dto
	createdLabel := Dto.ManifestLabel{
		Id:         createdLabelModel.UUID,
		Key:        createdLabelModel.Key,
		Value:      createdLabelModel.Value,
		SourceType: Common.MapLabelSourceType(createdLabelModel.SourceTypeId).Name,
		MediaType:  Common.MapMediaTypeName(createdLabelModel.MediaTypeId),
	}

	return createdLabel, nil
}
