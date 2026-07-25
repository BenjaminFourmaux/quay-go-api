package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func tagController() {
	registerRepositorySubRoute(http.MethodGet, "tag", listRepositoryTags)
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
	c.JSON(200, repository)
}
