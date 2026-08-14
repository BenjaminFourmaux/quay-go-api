package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func robotController() {
	userRobot := engine.Group("/api/v1/user/robots")
	{
		userRobot.Use(authorizedMiddleware)
		userRobot.GET("", listUserRobots)
	}

	organizationRobot := engine.Group("/api/v1/organization/:orgname/robots")
	{
		organizationRobot.Use(authorizedMiddleware)
		/*organizationRobot.GET("", listOrganizationRobots)*/
	}
}

// listUserRobots List user's robots accounts
// @Description List user's robots accounts
// @Summary List user's robots accounts
// @Tags Robot
// @Param token query bool false "Show robot token (default: false)"
// @Param repositories query bool false "Show robot repositories (default: false)"
// @Success 200 {object} []Dto.Robot
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots [get]
func listUserRobots(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	// Get filters from query params
	filters := extractFilters(c)

	organizations, err := Services.ListUserRobots(filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, organizations)
}
