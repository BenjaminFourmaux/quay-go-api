package Api

import (
	"github.com/gin-gonic/gin"
	dto "quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func usersController() {
	users := engine.Group("/api/v1/users")
	{
		users.Use(authorizedMiddleware)
		users.GET("/me", getCurrentUser)
	}

	// Just to avoid cleanup dependency
	_ = dto.UserMeResponse{}
}

// getCurrentUser Get the current authenticated user information
// @Description Get the current authenticated user information
// @Summary Get the current authenticated user information
// @Tags Users
// @Success 200 {object} dto.UserMeResponse
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/users/me [get]
func getCurrentUser(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{Auth.ReadUser})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	userId, _ := c.Get("authenticatedUserId")
	userScopes := Auth.ConvertListIdToScopes(c.GetString("scopes"))

	userInfo, err := Services.GetMeInfo(userId.(int), userScopes)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, userInfo)
}
