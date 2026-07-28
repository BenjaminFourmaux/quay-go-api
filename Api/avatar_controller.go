package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common/Errors"
	"quay-go-api/Services/Avatar"
	"strings"
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
// @Param format query string false "Format of the output (html, json, png). Default is html"
// @Success 200 {object} Dto.Avatar
// @Failure 400 {object} Errors.ErrorResponse "Bad Request invalid avatar kind"
// @Router /avatar/{kind}/{name} [get]
func getAvatar(c *gin.Context) {
	kind := c.Param("kind")
	name := c.Param("name")
	format := strings.ToLower(c.DefaultQuery("format", "html"))

	if kind != "user" && kind != "org" && kind != "team" {
		throwError(c, Errors.BadRequest("Invalid avatar kind. Must be 'user', 'org', or 'team'"))
		return
	}
	if format != "html" && format != "json" && format != "png" {
		throwError(c, Errors.BadRequest("Invalid format. Must be 'html', 'json', or 'png'"))
		return
	}

	// Get Avatar in the correct output format
	if format == "json" {
		avatar := Avatar.ToJSON(name, name, kind)
		c.JSON(http.StatusOK, avatar)
	} else if format == "png" {
		avatar, err := Avatar.ToPNG(name, name, 32, kind)
		if err != nil {
			throwError(c, Errors.InternalServerErrorWithMsg("Failed to generate PNG avatar"))
			return
		}
		c.Data(http.StatusOK, "image/png", avatar)
	} else {
		avatar := Avatar.ToHTML(name, name, 16, kind)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(avatar))
	}

}
