package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func robotController() {
	userRobot := engine.Group("/api/v1/user/robots")
	{
		userRobot.Use(authorizedMiddleware)
		userRobot.GET("", listUserRobots)
		userRobot.POST("", createUserRobot)
		userRobot.GET("/:robot_shortname", getUserRobot)
		userRobot.DELETE("/:robot_shortname", deleteUserRobot)
	}

	organizationRobot := engine.Group("/api/v1/organization/:orgname/robots")
	{
		organizationRobot.Use(authorizedMiddleware)
		/*organizationRobot.GET("", listOrganizationRobots)*/
	}
}

// listUserRobots List current user's robots accounts
// @Description List current user's robots accounts
// @Summary List current user's robots accounts
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

	robots, err := Services.ListUserRobots(filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, robots)
}

// createUserRobot Create for the current user a robot account
// @Description Create for the current user a robot account
// @Summary Create for the current user a robot account
// @Tags Robot
// @Accept json
// @Param message body Dto.CreateRobot true "Robot metadata"
// @Success 201 {object} Dto.Robot
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots [post]
func createUserRobot(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	var robotToCreate Dto.CreateRobot
	_ = c.BindJSON(&robotToCreate)

	newRobot, err := Services.CreateUserRobot(robotToCreate, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newRobot)
}

// getUserRobot Get a user robot
// @Description Get a user robot
// @Summary Get a user robot
// @Tags Robot
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 200 {object} Dto.Robot
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots/{robot_shortname} [get]
func getUserRobot(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	robotShortname := c.Param("robot_shortname")

	robot, err := Services.GetUserRobot(robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, robot)
}

// deleteUserRobot Delete a user robot
// @Description Delete a user robot
// @Summary Delete a user robot
// @Tags Robot
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 204
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots/{robot_shortname} [delete]
func deleteUserRobot(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	robotShortname := c.Param("robot_shortname")

	err := Services.DeleteUserRobot(robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
