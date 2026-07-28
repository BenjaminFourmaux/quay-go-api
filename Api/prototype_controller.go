package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func prototypeController() {
	prototype := engine.Group("/api/v1/organization/:orgname/prototypes")
	{
		prototype.Use(authorizedMiddleware)
		prototype.GET("", listPrototypes)
		prototype.POST("", createPrototype)
		prototype.GET("/:prototypeId", getPrototypeDetails)
		prototype.DELETE("/:prototypeId", deletePrototype)
		prototype.PATCH("/:prototypeId", updatePrototype)
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

// createPrototype Create a new permission prototype
// @Description Create a new permission prototype
// @Summary Create a new permission prototype
// @Tags Prototype
// @Accept json
// @Param orgname path string true "Name of the organization"
// @Param message body Dto.CreatePrototype true "Prototype metadata"
// @Success 201 {object} Dto.Prototype
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 404 {object} Errors.ErrorResponse "Not Found"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/prototypes [post]
func createPrototype(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	var prototypeToCreate Dto.CreatePrototype
	_ = c.BindJSON(&prototypeToCreate)

	newPrototype, err := Services.CreatePrototype(orgname, prototypeToCreate, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newPrototype)
}

// getPrototypeDetails Get a repository prototype by its ID
// @Description Get a repository prototype by its ID
// @Summary Get a repository prototype by its ID
// @Tags Prototype
// @Param orgname path string true "Name of the organization"
// @Param prototypeId path string true "ID of the prototype (UUID)"
// @Success 200 {object} Dto.Prototype
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/prototypes/{prototypeId} [get]
func getPrototypeDetails(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	prototypeId := c.Param("prototypeId")

	prototype, err := Services.GetPrototype(orgname, prototypeId, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, prototype)
}

// deletePrototype Delete a repository prototype by its ID
// @Description Delete a repository prototype by its ID
// @Summary Delete a repository prototype by its ID
// @Tags Prototype
// @Param orgname path string true "Name of the organization"
// @Param prototypeId path string true "ID of the prototype (UUID)"
// @Success 204
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/prototypes/{prototypeId} [delete]
func deletePrototype(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	prototypeId := c.Param("prototypeId")

	err := Services.DeletePrototype(orgname, prototypeId, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// updatePrototype Update permission prototype role
// @Description Update permission prototype role
// @Summary  Update permission prototype role
// @Tags Prototype
// @Accept json
// @Param orgname path string true "Name of the organization"
// @Param prototypeId path string true "ID of the prototype (UUID)"
// @Param update body Dto.UpdatePrototype true "Prototype details to change"
// @Success 200 {object} Dto.Prototype
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/prototypes/{prototypeId} [patch]
func updatePrototype(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.OrgAdmin})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")
	prototypeId := c.Param("prototypeId")

	var prototypeToUpdate Dto.UpdatePrototype
	_ = c.BindJSON(&prototypeToUpdate)

	updatedPrototype, err := Services.UpdatePrototype(orgname, prototypeId, prototypeToUpdate, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedPrototype)
}
