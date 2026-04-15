package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func messagesController() {
	messages := engine.Group("/api/v1/messages")
	{
		messages.Use(authorizedMiddleware)
		messages.GET("/", listMessages)
		/*messages.POST("/", createMessage)
		messages.PATCH("/:id", updateMessage)
		messages.DELETE("/:id", deleteMessage)*/
	}

	// Just to avoid cleanup dependency
	_ = dto.Message{}
}

// listMessages List messages displayed in the quay web app for all user
// @Description List messages displayed in the quay web app for all user
// @Summary List messages displayed in the quay web app for all user
// @Tags Messages
// @Success 200 {object} []dto.Message
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/messages [get]
func listMessages(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{Auth.AdminUser})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	messages, err := Services.ListMessages()
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}
