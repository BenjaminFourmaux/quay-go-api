package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common/Errors"
	"quay-go-api/Services/Avatar"
)

func avatarController() {
	avatar := engine.Group("/avatar")
	{
		avatar.GET("/:kind/:name", getAvatar)
	}
}

// getAvatar Get user/org/team avatar
// @Summary Get user/org/team avatar
// @Description Get user/org/team avatar
// @Tags Avatar
// @Param kind path string true "Avatar kind (user, org, team)"
// @Param name path string true "Username/Org name/Team name to get avatar for"
// @Success 200 {html} avatar
// @Failure 400 {object} Errors.ErrorResponse "Bad Request invalid avatar kind"
// @Router /avatar/{kind}/{name} [get]
func getAvatar(c *gin.Context) {
	kind := c.Param("kind")
	name := c.Param("name")

	if kind != "user" && kind != "org" && kind != "team" {
		throwError(c, Errors.BadRequest("Invalid avatar kind. Must be 'user', 'org', or 'team'"))
		return
	}

	avatar := Avatar.GetHTML(name, name, 16, kind)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(avatar))
}
