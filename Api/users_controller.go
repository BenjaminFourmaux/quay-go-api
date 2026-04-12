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
}

// getCurrentUser Get the current authenticated user information
// @Description Get the current authenticated user information
// @Summary Get the current authenticated user information
// @Tags Users
// @Success 200 {object} dto.UserMeResponse
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/users/me [get]
func getCurrentUser(c *gin.Context) {
	_ = dto.UserMeResponse{}
	hasScopes := requiredScopes(c, []Auth.Scope{Auth.ReadUser})
	if hasScopes != nil {
		throwError(c, hasScopes)
		return
	}

	userId, _ := c.Get("authenticatedUserId")

	userInfo, err := Services.GetMeInfo(userId.(int))
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, userInfo)
}
