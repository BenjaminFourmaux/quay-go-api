package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func organizationController() {
	organization := engine.Group("/api/v1/organization")
	{
		organization.Use(authorizedMiddleware)
		organization.GET("/", listOrganizations)
		organization.POST("/", createOrganization)
		organization.GET("/:orgname", getOrganizationDetails)
	}
}

// listOrganizations List user's organizations
// @Description List user's organizations
// @Summary List user's organizations
// @Tags Organization
// @Success 200 {object} []Dto.UserOrganization
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization [get]
func listOrganizations(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	userId, _ := c.Get("authenticatedUserId")
	userScopesInterface, _ := c.Get("scopes")
	userScopes := Common.ConvertScopeStringInAuthScopes(userScopesInterface.(string))

	organizations, err := Services.GetUserOrganizations(userId.(int), userScopes)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, organizations)
}

// createOrganization Create a new organization
// @Description Create a new organization
// @Summary Create a new organization
// @Tags Organization
// @Accept json
// @Param message body Dto.CreateOrganization true "Organization metadata"
// @Success 201 {object} Dto.Organization
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization [post]
func createOrganization(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	userId, _ := c.Get("authenticatedUserId")
	userScopesInterface, _ := c.Get("scopes")
	userScopes := Common.ConvertScopeStringInAuthScopes(userScopesInterface.(string))

	var organizationToCreate Dto.CreateOrganization
	_ = c.BindJSON(&organizationToCreate)

	newOrganization, err := Services.CreateOrganization(organizationToCreate, userId.(int), userScopes)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newOrganization)
}

// getOrganizationDetails Get details in an organization
// @Description Get details in an organization
// @Summary Get details in an organization
// @Tags Organization
// @Param orgname path string true "Name of the organization"
// @Success 200 {object} Dto.Organization
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname} [get]
func getOrganizationDetails(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	userId, _ := c.Get("authenticatedUserId")
	userScopesInterface, _ := c.Get("scopes")
	userScopes := Common.ConvertScopeStringInAuthScopes(userScopesInterface.(string))

	orgname := c.Param("orgname")

	organization, err := Services.GetOrganizationDetailsByName(orgname, userId.(int), userScopes)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, organization)
}
