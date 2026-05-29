package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func membersController() {
	teams := engine.Group("/api/v1/organization/:orgname/team/:teamname/members")
	{
		teams.Use(authorizedMiddleware)
		teams.GET("/", listTeamMembers)
		/*teams.POST("/", createOrganizationTeam)
		teams.GET("/:teamname", getOrganizationTeam)
		teams.PATCH("/:teamname", updateOrganizationTeam)
		teams.DELETE("/:teamname", deleteOrganizationTeam)*/
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
