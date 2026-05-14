package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func usersController() {
	users := engine.Group("/api/v1/users")
	{
		users.Use(authorizedMiddleware)
		users.GET("/me", getCurrentUser)
	}
}

// getCurrentUser Get the current authenticated user information
// @Description Get the current authenticated user information
// @Summary Get the current authenticated user information
// @Tags Users
// @Success 200 {object} Dto.UserMeResponse
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/users/me [get]
func getCurrentUser(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	userInfo, err := Services.GetMeInfo(currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, userInfo)
}
