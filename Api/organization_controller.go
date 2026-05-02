package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common/Errors"
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
		organization.DELETE("/:orgname", deleteOrganization)
		organization.PATCH("/:orgname", updateOrganization)
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
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	organizations, err := Services.GetUserOrganizations(currentUser)
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
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	var organizationToCreate Dto.CreateOrganization
	_ = c.BindJSON(&organizationToCreate)

	newOrganization, err := Services.CreateOrganization(organizationToCreate, currentUser)
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
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	organization, err := Services.GetOrganizationDetailsByName(orgname, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, organization)
}

// deleteOrganization Delete an organization
// @Description Delete an organization
// @Summary Delete an organization
// @Tags Organization
// @Param orgname path string true "Name of the organization"
// @Success 204 "No Content"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname} [delete]
func deleteOrganization(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	err := Services.DeleteOrganization(orgname, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// updateOrganization Update organization details
// @Description Update organization details
// @Summary Update organization details
// @Tags Organization
// @Accept json
// @Param orgname path string true "Name of the organization"
// @Param message body Dto.UpdateOrganization true "Organization details to change"
// @Success 200 {object} Dto.Organization
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname} [patch]
func updateOrganization(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	var updateOrganization Dto.UpdateOrganization
	if err := c.BindJSON(&updateOrganization); err != nil {
		throwError(c, Errors.RequestBodyInvalid())
		return
	}

	updatedOrganization, err := Services.UpdateOrganization(orgname, updateOrganization, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedOrganization)
}
