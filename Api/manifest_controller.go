package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func manifestController() {
	registerRepositorySubRoute(http.MethodGet, "manifest/:manifestRef", getManifest)
	registerRepositorySubRoute(http.MethodGet, "manifest/:manifestRef/labels", getManifestLabels)
	registerRepositorySubRoute(http.MethodPost, "manifest/:manifestRef/labels", createManifestLabel)
}

// getManifest Get a repository manifest
// @Description Get a repository manifest
// @Summary Get a repository manifest
// @Tags Manifest
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param manifestRef path string true "Manifest reference sha256"
// @Success 200 {object} Dto.Manifest
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/manifest/{manifestRef} [get]
func getManifest(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	manifestRef := c.Param("manifestRef")

	repository, err := Services.GetManifest(repositoryNamespaced, manifestRef, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, repository)
}

// getManifestLabels Get a repository manifest labels
// @Description Get a repository manifest labels
// @Summary Get a repository manifest labels
// @Tags Manifest
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param manifestRef path string true "Manifest reference sha256"
// @Success 200 {object} Dto.Manifest
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/manifest/{manifestRef}/labels [get]
func getManifestLabels(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	manifestRef := c.Param("manifestRef")

	repository, err := Services.GetManifestLabels(repositoryNamespaced, manifestRef, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, repository)
}

// createManifestLabel Create a new manifest label
// @Description Create a new manifest label
// @Summary Create a new manifest label
// @Tags Manifest
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param manifestRef path string true "Manifest reference sha256"
// @Param label body Dto.AddManifestLabel true "Manifest label data"
// @Success 201 {object} Dto.AddManifestLabel
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/manifest/{manifestRef}/labels [post]
func createManifestLabel(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	manifestRef := c.Param("manifestRef")

	var addLabel Dto.AddManifestLabel
	_ = c.BindJSON(&addLabel)

	createdLabel, err := Services.CreateManifestLabel(repositoryNamespaced, manifestRef, addLabel, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(201, createdLabel)
}
