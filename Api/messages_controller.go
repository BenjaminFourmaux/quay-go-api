package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func messagesController() {
	messages := engine.Group("/api/v1/messages")
	{
		messages.Use(authorizedMiddleware)
		messages.GET("/", listMessages)
		messages.POST("/", createMessage)
		/*messages.PATCH("/:uuid", updateMessage)
		messages.DELETE("/:uuid", deleteMessage)*/
	}
}

// listMessages List messages displayed on the quay web app for all user
// @Description List messages displayed on the quay web app for all user
// @Summary List messages displayed on the quay web app for all user
// @Tags Messages
// @Success 200 {object} []Dto.Message
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

// createMessage Create a new message to displayed on the quay web app for all user
// @Description Create a new message to displayed on the quay web app for all user
// @Summary Create a new message to displayed on the quay web app for all user
// @Tags Messages
// @Accept json
// @Param message body Dto.CreateMessage true "Message content and severity"
// @Success 201 {object} Dto.Message
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/messages [post]
func createMessage(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{Auth.AdminUser})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	var createMessage Dto.CreateMessage
	_ = c.BindJSON(&createMessage)

	createdMessage, err := Services.CreateMessage(createMessage)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusCreated, createdMessage)
}
