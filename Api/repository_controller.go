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

// repositorySubRoutePattern describes a sub-resource path nested under a
// repository (e.g. "permissions/team/:team"), split into segments so it can
// be matched against the trailing segments of the captured wildcard path.
type repositorySubRoutePattern struct {
	segments []string
	handler  func(c *gin.Context, repositoryNamespaced string)
}

// repositorySubRoutes maps an HTTP method to the sub-resource patterns
// registered for it. Gin only allows a single trailing "*wildcard" per HTTP
// method under a given path, so sub-resources (permissions, and any future
// ones) register themselves here instead of declaring their own conflicting
// wildcard route.
var repositorySubRoutes = map[string][]repositorySubRoutePattern{}

// registerRepositorySubRoute lets other controllers (e.g. permission_controller.go)
// plug additional endpoints under /api/v1/repository/{repository}/... for any
// HTTP method, without creating a conflicting Gin wildcard route.
//
// pattern uses the same ":name" syntax as Gin for trailing dynamic segments,
// e.g. "permissions/team" or "permissions/team/:team". Captured values are
// available via c.Param(name) exactly like a normal Gin route.
func registerRepositorySubRoute(method string, pattern string, handler func(c *gin.Context, repositoryNamespaced string)) {
	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	repositorySubRoutes[method] = append(repositorySubRoutes[method], repositorySubRoutePattern{
		segments: segments,
		handler:  handler,
	})
}

func repositoryController() {
	repository := engine.Group("/api/v1/repository")
	{
		repository.Use(authorizedMiddleware)
		repository.GET("", listRepositories)
		repository.POST("", createRepository)
		repository.Any("/*repository", dispatchRepository)
	}
}

// dispatchRepository inspects the wildcard path captured after /repository/
// and routes it to a registered sub-route handler (e.g. permissions) when its
// trailing segments match, falling back to plain repository CRUD otherwise.
// This allows both "myglobalrepo/permissions/team/myteam" and
// "myorg/myrepo/permissions/team/myteam" to be handled - for any HTTP verb -
// without needing one route per repository name shape.
func dispatchRepository(c *gin.Context) {
	path := strings.Trim(c.Param("repository"), "/")
	pathSegments := strings.Split(path, "/")

	if dispatchToSubRoute(c, pathSegments) {
		return
	}

	switch c.Request.Method {
	case http.MethodGet:
		getRepository(c, path)
	case http.MethodPatch:
		updateRepository(c, path)
	case http.MethodDelete:
		deleteRepository(c, path)
	default:
		c.Status(http.StatusMethodNotAllowed)
	}
}

// dispatchToSubRoute tries every possible split point between the repository
// name and a registered sub-route's segments (since the repository name
// itself may be 1 segment - "myrepo" - or 2 - "myorg/myrepo"). It returns true
// once a pattern matching the current HTTP method has consumed the trailing
// segments.
func dispatchToSubRoute(c *gin.Context, pathSegments []string) bool {
	for _, pattern := range repositorySubRoutes[c.Request.Method] {
		if len(pattern.segments) >= len(pathSegments) {
			continue
		}

		splitAt := len(pathSegments) - len(pattern.segments)
		if !matchSubRouteSegments(c, pattern.segments, pathSegments[splitAt:]) {
			continue
		}

		repositoryNamespaced := strings.Join(pathSegments[:splitAt], "/")
		pattern.handler(c, repositoryNamespaced)
		return true
	}

	return false
}

// matchSubRouteSegments compares pattern segments against the actual trailing
// path segments, capturing ":name" segments as Gin params along the way.
func matchSubRouteSegments(c *gin.Context, pattern []string, actual []string) bool {
	for i, segment := range pattern {
		if strings.HasPrefix(segment, ":") {
			c.Params = append(c.Params, gin.Param{Key: segment[1:], Value: actual[i]})
			continue
		}
		if segment != actual[i] {
			return false
		}
	}
	return true
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
func updateRepository(c *gin.Context, repositoryName string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.AdminRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

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
func deleteRepository(c *gin.Context, repositoryName string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.AdminRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	err := Services.DeleteRepository(repositoryName, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
