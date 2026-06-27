package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func repositoryController() {
	repository := engine.Group("/api/v1/repository")
	{
		repository.Use(authorizedMiddleware)
		repository.GET("", listRepositories)
		repository.POST("", createRepository)
		repository.GET("/*repository", getRepository)
	}
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
func getRepository(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	repositoryNamespaced := strings.TrimPrefix(c.Param("repository"), "/")

	// Get filters from query params
	filters := extractFilters(c)

	repository, err := Services.GetRepository(repositoryNamespaced, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, repository)
}
