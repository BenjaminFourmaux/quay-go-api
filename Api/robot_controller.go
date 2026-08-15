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
		userRobot.GET("/:robot_shortname/permissions", listUserRobotPermissions)
		userRobotFederation := userRobot.Group("/:robot_shortname/federation")
		{
			userRobotFederation.Use(authorizedMiddleware)
			userRobotFederation.GET("", listUserRobotFederation)
			userRobotFederation.POST("", createOrUpdateUserRobotFederation)
			userRobotFederation.DELETE("", deleteUserRobotFederation)
		}
	}

	organizationRobot := engine.Group("/api/v1/organization/:orgname/robots")
	{
		organizationRobot.Use(authorizedMiddleware)
		organizationRobot.GET("", listOrganizationRobots)
		organizationRobot.POST("", createOrganizationRobot)
		organizationRobot.GET("/:robot_shortname", getOrganizationRobot)
		organizationRobot.DELETE("/:robot_shortname", deleteOrganizationRobot)
		organizationRobot.GET("/:robot_shortname/permissions", listOrganizationRobotPermissions)
		organizationRobotFederation := organizationRobot.Group("/:robot_shortname/federation")
		{
			organizationRobotFederation.Use(authorizedMiddleware)
			organizationRobotFederation.GET("", listOrganizationRobotFederation)
			organizationRobotFederation.POST("", createOrUpdateOrganizationRobotFederation)
			organizationRobotFederation.DELETE("", deleteOrganizationRobotFederation)
		}
	}
}

// <editor-fold desc="User (current user)">

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

// listUserRobotPermissions List a user robot permissions
// @Description List a user robot permissions
// @Summary List a user robot permissions
// @Tags Robot
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 200 {object} []Dto.RobotPermission
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots/{robot_shortname}/permissions [get]
func listUserRobotPermissions(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	robotShortname := c.Param("robot_shortname")

	newRobot, err := Services.GetUserRobotPermissions(robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, newRobot)
}

// listUserRobotFederation List a user robot federation
// @Description List a user robot federation
// @Summary List a user robot federation
// @Tags Robot
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 200 {object} []Dto.RobotFederation
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots/{robot_shortname}/federation [get]
func listUserRobotFederation(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	robotShortname := c.Param("robot_shortname")

	permissions, err := Services.ListUserRobotFederations(robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// createOrUpdateUserRobotFederation Create or update a user robot account federations
// @Description Create or update a user robot account federations
// @Summary Create or update a user robot account federations
// @Tags Robot
// @Accept json
// @Param robot_shortname path string true "Shortname of the robot"
// @Param message body []Dto.RobotFederation true "Federations metadata"
// @Success 201 {object} []Dto.RobotFederation
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots/{robot_shortname}/federation [post]
func createOrUpdateUserRobotFederation(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	robotShortname := c.Param("robot_shortname")

	var robotFederationsToCreate []Dto.RobotFederation
	_ = c.BindJSON(&robotFederationsToCreate)

	newFederations, err := Services.CreateOrUpdateUserRobotFederations(robotShortname, robotFederationsToCreate, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newFederations)
}

// deleteUserRobotFederation Delete a user robot federation
// @Description Delete a user robot federation
// @Summary Delete a user robot federation
// @Tags Robot
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 204
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/user/robots/{robot_shortname}/federation [delete]
func deleteUserRobotFederation(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	robotShortname := c.Param("robot_shortname")

	err := Services.DeleteUserRobotFederation(robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// </editor-fold>

// <editor-fold desc="Organization">

// listOrganizationRobots List organization's robots accounts
// @Description List organization's robots accounts
// @Summary List organization's robots accounts
// @Tags Robot
// @Param orgname path string true "Name of the organization"
// @Param token query bool false "Show robot token (default: false)"
// @Param repositories query bool false "Show robot repositories (default: false)"
// @Success 200 {object} []Dto.Robot
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots [get]
func listOrganizationRobots(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")

	// Get filters from query params
	filters := extractFilters(c)

	robots, err := Services.ListOrganizationRobots(orgName, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, robots)
}

// createOrganizationRobot Create for the organization a robot account
// @Description Create for the organization a robot account
// @Summary Create for the organization a robot account
// @Tags Robot
// @Accept json
// @Param orgname path string true "Name of the organization"
// @Param message body Dto.CreateRobot true "Robot metadata"
// @Success 201 {object} Dto.Robot
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots [post]
func createOrganizationRobot(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	var robotToCreate Dto.CreateRobot
	_ = c.BindJSON(&robotToCreate)

	orgName := c.Param("orgname")

	newRobot, err := Services.CreateOrganizationRobot(orgName, robotToCreate, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newRobot)
}

// getOrganizationRobot Get a organization robot account
// @Description Get a organization robot account
// @Summary Get a organization robot account
// @Tags Robot
// @Param orgname path string true "Name of the organization"
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 200 {object} Dto.Robot
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots/{robot_shortname} [get]
func getOrganizationRobot(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")
	robotShortname := c.Param("robot_shortname")

	robot, err := Services.GetOrganizationRobot(orgName, robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, robot)
}

// deleteOrganizationRobot Delete a organization robot
// @Description Delete a organization robot
// @Summary Delete a organization robot
// @Tags Robot
// @Param orgname path string true "Name of the organization"
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 204
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots/{robot_shortname} [delete]
func deleteOrganizationRobot(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")
	robotShortname := c.Param("robot_shortname")

	err := Services.DeleteOrganizationRobot(orgName, robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// listOrganizationRobotPermissions List a organization robot permissions
// @Description List a organization robot permissions
// @Summary List a organization robot permissions
// @Tags Robot
// @Param orgname path string true "Name of the organization"
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 200 {object} []Dto.RobotPermission
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots/{robot_shortname}/permissions [get]
func listOrganizationRobotPermissions(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")
	robotShortname := c.Param("robot_shortname")

	newRobot, err := Services.GetOrganizationRobotPermissions(orgName, robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, newRobot)
}

// listOrganizationRobotFederation List an organization robot federation
// @Description List an organization robot federation
// @Summary List an organization robot federation
// @Tags Robot
// @Param orgname path string true "Name of the organization"
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 200 {object} []Dto.RobotFederation
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots/{robot_shortname}/federation [get]
func listOrganizationRobotFederation(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")
	robotShortname := c.Param("robot_shortname")

	permissions, err := Services.ListOrganizationRobotFederations(orgName, robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// createOrUpdateOrganizationRobotFederation Create or update an organization robot account federations
// @Description Create or update an organization robot account federations
// @Summary Create or update an organization robot account federations
// @Tags Robot
// @Accept json
// @Param orgname path string true "Name of the organization"
// @Param robot_shortname path string true "Shortname of the robot"
// @Param message body []Dto.RobotFederation true "Federations metadata"
// @Success 201 {object} []Dto.RobotFederation
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots/{robot_shortname}/federation [post]
func createOrUpdateOrganizationRobotFederation(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")
	robotShortname := c.Param("robot_shortname")

	var robotFederationsToCreate []Dto.RobotFederation
	_ = c.BindJSON(&robotFederationsToCreate)

	newFederations, err := Services.CreateOrUpdateOrganizationRobotFederations(orgName, robotShortname, robotFederationsToCreate, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newFederations)
}

// deleteOrganizationRobotFederation Delete an organization robot federation
// @Description Delete an organization robot federation
// @Summary Delete an organization robot federation
// @Tags Robot
// @Param orgname path string true "Name of the organization"
// @Param robot_shortname path string true "Shortname of the robot"
// @Success 204
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/robots/{robot_shortname}/federation [delete]
func deleteOrganizationRobotFederation(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgName := c.Param("orgname")
	robotShortname := c.Param("robot_shortname")

	err := Services.DeleteOrganizationRobotFederation(orgName, robotShortname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// </editor-fold>
