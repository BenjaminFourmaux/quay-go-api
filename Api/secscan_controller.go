package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func secScanController() {
	registerRepositorySubRoute(http.MethodGet, "manifest/:manifestRef/security", getRepositoryManifestSecScan)

}

// getRepositoryManifestSecScan Get security scan report for a manifest
// @Description Get security scan report for a manifest
// @Summary Get security scan report for a manifest
// @Tags SecurityScan
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param manifestRef path string true "Manifest reference sha256"
// @Param include_cve query bool false "Include CVEs information"
// @Success 200 {object} Dto.SecScanReport
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/manifest/{manifestRef}/security [get]
func getRepositoryManifestSecScan(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	manifestRef := c.Param("manifestRef")

	// Extract filters from query params
	filters := extractFilters(c)

	repository, err := Services.GetRepositoryManifestSecScan(repositoryNamespaced, manifestRef, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, repository)
}
