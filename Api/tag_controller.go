package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func tagController() {
	registerRepositorySubRoute(http.MethodGet, "tag", listRepositoryTags)
	registerRepositorySubRoute(http.MethodGet, "tag/:tag", getRepositoryTag)
	registerRepositorySubRoute(http.MethodPut, "tag/:tag", updateRepositoryTag)
	registerRepositorySubRoute(http.MethodDelete, "tag/:tag", deleteRepositoryTag)
}

// listRepositoryTags List tags on a repository
// @Description List tags on a repository
// @Summary List tags on a repository
// @Tags Tag
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param include_vulnerabilities query bool false "Include vulnerability information"
// @Success 200 {object} []Dto.Tag
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/tag [get]
func listRepositoryTags(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	// Extract filters from query params
	filters := extractFilters(c)

	repository, err := Services.ListRepositoryTags(repositoryNamespaced, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, repository)
}

// getRepositoryTag Get a specific tag from a repository
// @Description Get a specific tag from a repository
// @Summary Get a specific tag from a repository
// @Tags Tag
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param tag path string true "Tag name"
// @Param include_vulnerabilities query bool false "Include vulnerability information"
// @Success 200 {object} Dto.Tag
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/tag/{tag} [get]
func getRepositoryTag(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	tag := c.Param("tag")

	// Extract filters from query params
	filters := extractFilters(c)

	repository, err := Services.GetRepositoryTag(repositoryNamespaced, tag, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, repository)
}

// updateRepositoryTag Update a specific tag in a repository
// @Description Update a specific tag in a repository
// @Summary Update a specific tag in a repository
// @Tags Tag
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param tag path string true "Tag name"
// @Param update body Dto.UpdateTag true "Tag details to change"
// @Success 200 {object} Dto.Tag
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/tag/{tag} [put]
func updateRepositoryTag(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	tag := c.Param("tag")

	var updateTag Dto.UpdateTag
	if err := c.BindJSON(&updateTag); err != nil {
		throwError(c, Errors.RequestBodyInvalid())
		return
	}

	repository, err := Services.UpdateRepositoryTag(repositoryNamespaced, tag, updateTag, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, repository)
}

// deleteRepositoryTag Delete a specific tag in a repository
// @Description Delete a specific tag in a repository
// @Summary Delete a specific tag in a repository
// @Tags Tag
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param tag path string true "Tag name"
// @Success 204
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/tag/{tag} [delete]
func deleteRepositoryTag(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	tag := c.Param("tag")

	err := Services.DeleteRepositoryTag(repositoryNamespaced, tag, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
