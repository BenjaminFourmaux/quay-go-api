package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
	"strings"
)

// repositorySubRouteHandlers maps a URL suffix (e.g. "/permissions/team") to a
// handler receiving the parsed repository name (namespace/repository or just
// repository). Gin only allows a single "*wildcard" per HTTP method on a given
// path, so sub-resources under /repository/... (like permissions) register
// themselves here instead of declaring their own conflicting wildcard route.
var repositorySubRouteHandlers = map[string]func(c *gin.Context, repositoryNamespaced string){}

// registerRepositorySubRoute lets other controllers (e.g. permission_controller.go)
// plug additional GET endpoints under /api/v1/repository/{repository}/... without
// creating a conflicting Gin wildcard route.
func registerRepositorySubRoute(suffix string, handler func(c *gin.Context, repositoryNamespaced string)) {
	repositorySubRouteHandlers[suffix] = handler
}

func repositoryController() {
	repository := engine.Group("/api/v1/repository")
	{
		repository.Use(authorizedMiddleware)
		repository.GET("", listRepositories)
		repository.POST("", createRepository)
		repository.GET("/*repository", dispatchRepositoryGet)
		repository.PATCH("/*repository", updateRepository)
		repository.DELETE("/*repository", deleteRepository)
	}
}

// dispatchRepositoryGet inspects the wildcard path captured after /repository/
// and routes it to a registered sub-route handler (e.g. permissions) when its
// suffix matches, falling back to plain repository retrieval otherwise. This
// allows both "myglobalrepo/permissions/team" and "myorg/myrepo/permissions/team"
// to be handled without needing one route per repository name shape.
func dispatchRepositoryGet(c *gin.Context) {
	path := strings.TrimPrefix(c.Param("repository"), "/")

	for suffix, handler := range repositorySubRouteHandlers {
		if repositoryNamespaced, ok := strings.CutSuffix(path, suffix); ok {
			handler(c, repositoryNamespaced)
			return
		}
	}

	getRepository(c, path)
}

// listRepositories List repositories
// @Description List repositories
// @Summary List repositories
// @Tags Repository
// @Param namespace query string false "Organization or Username as namespace to filter repositories"
// @Param is_public query bool false "Filter on public repositories"
// @Param is_starred query bool false "Filter on stared repositories"
// @Success 200 {object} []Dto.Repository
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository [get]
func listRepositories(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	// Get filters from query params
	filters := extractFilters(c)

	listTeams, err := Services.ListRepositories(filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, listTeams)
}

// createRepository Create a repository
// @Description Create a repository
// @Summary Create a repository
// @Tags Repository
// @Accept json
// @Param message body Dto.CreateRepository true "Repository metadata"
// @Success 201 {object} Dto.Repository
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository [post]
func createRepository(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.CreateRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	var repositoryToCreate Dto.CreateRepository
	_ = c.BindJSON(&repositoryToCreate)

	newRepository, err := Services.CreateRepository(repositoryToCreate, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newRepository)
}

// getRepository Get a repository
// @Description Get a repository
// @Summary Get a repository
// @Tags Repository
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param include_tags query bool false "Include tags in the repository response"
// @Param include_stats query bool false "Include statistics in the repository response"
// @Success 200 {object} Dto.RepositoryDetails
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository} [get]
func getRepository(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	// Get filters from query params
	filters := extractFilters(c)

	repository, err := Services.GetRepository(repositoryNamespaced, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, repository)
}

// updateRepository Update repository details
// @Description Update repository details
// @Summary Update repository details
// @Tags Repository
// @Accept json
// @Param repository path string true "Name of the repository"
// @Param message body Dto.UpdateRepository true "Repository details to change"
// @Success 200 {object} Dto.Repository
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository} [patch]
func updateRepository(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.AdminRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	repositoryName := strings.TrimPrefix(c.Param("repository"), "/")

	var updateRepository Dto.UpdateRepository
	if err := c.BindJSON(&updateRepository); err != nil {
		throwError(c, Errors.RequestBodyInvalid())
		return
	}

	updatedRepository, err := Services.UpdateRepository(repositoryName, updateRepository, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedRepository)
}

// deleteRepository Delete a repository
// @Description Delete a repository
// @Summary Delete a repository
// @Tags Repository
// @Accept json
// @Param repository path string true "Name of the repository"
// @Success 204
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository} [delete]
func deleteRepository(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.AdminRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	repositoryName := strings.TrimPrefix(c.Param("repository"), "/")

	err := Services.DeleteRepository(repositoryName, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
