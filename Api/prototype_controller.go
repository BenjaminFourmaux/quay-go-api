package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func prototypeController() {
	prototype := engine.Group("/api/v1/organization/:orgname/prototypes")
	{
		prototype.Use(authorizedMiddleware)
		prototype.GET("", listPrototypes)
		/*prototype.POST("", createPrototype)
		prototype.GET("/:prototypeId", getPrototypeDetails)
		prototype.DELETE("/:prototypeId", deletePrototype)
		prototype.PATCH("/:prototypeId", updatePrototype)*/
	}
}

// listPrototypes List the existing prototypes for this organization
// @Description List the existing prototypes for this organization
// @Summary List the existing prototypes for this organization
// @Tags Prototype
// @Param orgname path string true "Name of the organization"
// @Success 200 {object} []Dto.Prototype
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/prototypes [get]
func listPrototypes(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	prototypes, err := Services.ListOrganizationPrototypes(orgname, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}

	c.JSON(200, prototypes)
}
