package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func teamController() {
	teams := engine.Group("/api/v1/organization/:orgname/team")
	{
		teams.Use(authorizedMiddleware)
		teams.GET("/", listOrganizationTeams)
		teams.POST("/", createOrganizationTeam)
		teams.PATCH("/:teamname", updateOrganizationTeam)
		teams.DELETE("/:teamname", deleteOrganizationTeam)
	}
}

// listOrganizationTeams List organization's teams
// @Description List organization's teams with optional filtering
// @Summary List organization's teams
// @Tags Team
// @Param orgname path string true "Name of the organization"
// @Param role query string false "Filter teams by role name (e.g., 'admin', 'creator', 'member')"
// @Param name query string false "Filter teams by name"
// @Success 200 {object} []Dto.Team
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team [get]
func listOrganizationTeams(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	// Get filters from query params
	filters := extractFilters(c)

	listTeams, err := Services.ListTeamsOfOrganization(orgname, filters, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, listTeams)
}

// createOrganizationTeam Create a team inside the organization
// @Description Create a team inside the organization
// @Summary Create a team inside the organization
// @Tags Team
// @Accept json
// @Param orgname path string true "Name of the organization"
// @Param message body Dto.CreateTeam true "Team to create"
// @Success 201 {object} Dto.Team
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team [post]
func createOrganizationTeam(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	var teamToCreate Dto.CreateTeam
	_ = c.BindJSON(&teamToCreate)

	newTeam, err := Services.CreateTeam(teamToCreate, orgname, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newTeam)
}

// updateOrganizationTeam Update team details
// @Description Update team details
// @Summary Update team details
// @Tags Team
// @Accept json
// @Param orgname path string true "Organization name"
// @Param teamname path string true "Team name to update"
// @Param team body Dto.UpdateTeam true "Team description and role to update"
// @Success 200 {object} Dto.Team
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team/{teamname} [patch]
func updateOrganizationTeam(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	teamname := c.Param("teamname")

	var teamToUpdate Dto.UpdateTeam
	_ = c.BindJSON(&teamToUpdate)

	updatedTeam, err := Services.UpdateTeam(teamToUpdate, orgname, teamname, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedTeam)
}

// deleteOrganizationTeam Delete a team
// @Description Delete a team
// @Summary Delete a team
// @Tags Team
// @Param orgname path string true "Name of the organization"
// @Param teamname path string true "Name of the team to delete"
// @Success 204 "No Content"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team/{teamname} [delete]
func deleteOrganizationTeam(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	teamname := c.Param("teamname")

	err := Services.DeleteTeam(orgname, teamname, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
