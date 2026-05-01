package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Common"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func organizationController() {
	organization := engine.Group("/api/v1/organization")
	{
		organization.Use(authorizedMiddleware)
		organization.GET("/", listOrganizations)
	}
}

// listOrganizations List user's organizations
// @Description List user's organizations
// @Summary List user's organizations
// @Tags Organization
// @Success 200 {object} []Dto.Organization
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization [get]
func listOrganizations(c *gin.Context) {
	hasScopesErr := requiredScopes(c, []Auth.Scope{})
	if hasScopesErr != nil {
		throwError(c, hasScopesErr)
		return
	}

	userId, _ := c.Get("authenticatedUserId")
	userScopesInterface, _ := c.Get("scopes")

	userScopes := Common.ConvertScopeStringInAuthScopes(userScopesInterface.(string))

	organization, err := Services.GetUserOrganization(userId.(int), userScopes)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, organization)
}
