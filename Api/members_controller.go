package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func membersController() {
	members := engine.Group("/api/v1/organization/:orgname/team/:teamname/members")
	{
		members.Use(authorizedMiddleware)
		members.GET("", listTeamMembers)
		members.PUT("/:membername", addMemberToTeam)
		members.DELETE("/:membername", removeMemberFromTeam)
	}
}

// listTeamMembers List team's members
// @Description List team's members
// @Summary List team's members
// @Tags Members
// @Param orgname path string true "Name of the organization"
// @Param teamname path string true "Name of the team"
// @Param is_robot query bool false "Filter team members on is_robot"
// @Param is_invited query bool false "Filter team members on invited status"
// @Success 200 {object} []Dto.TeamMember
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team/{teamname}/members [get]
func listTeamMembers(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	teamname := c.Param("teamname")

	// Get filters from query params
	filters := extractFilters(c)

	listMembers, err := Services.ListTeamMembers(orgname, teamname, filters, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, listMembers)
}

// addMemberToTeam Add a user to a team
// @Description Add a user to a team
// @Summary Add a user to a team
// @Tags Members
// @Param orgname path string true "Name of the organization"
// @Param teamname path string true "Name of the team"
// @Param membername path string true "Name of the user"
// @Success 201 {object} Dto.TeamMember
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team/{teamname}/members/{membername} [put]
func addMemberToTeam(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	teamname := c.Param("teamname")
	membername := c.Param("membername")

	newMember, err := Services.AddMemberToTeam(orgname, teamname, membername, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(201, newMember)
}

// removeMemberFromTeam Remove a user from a team
// @Description Remove a user from a team
// @Summary Remove a user from a team
// @Tags Members
// @Param orgname path string true "Name of the organization"
// @Param teamname path string true "Name of the team"
// @Param membername path string true "Name of the user"
// @Success 204 "No Content"
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/team/{teamname}/members/{membername} [delete]
func removeMemberFromTeam(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	teamname := c.Param("teamname")
	membername := c.Param("membername")

	err := Services.RemoveMemberToTeam(orgname, teamname, membername, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.Status(204)
}
