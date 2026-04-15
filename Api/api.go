package Api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"quay-go-api/Common/Errors"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Logger"
	"strings"
)

var engine *gin.Engine

func StartServer() {
	engine = gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	Logger.Info("Starting server on port " + port)

	endpointsRegistration()

	err := engine.Run(":" + port)
	if err != nil {
		Logger.Error("Failed to start server: " + err.Error())
		return
	}
}

// <editor-fold desc="Private functions">

func endpointsRegistration() {
	healthController()
	messagesController()
	usersController()

	// Add Swagger endpoint
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

/*
authorizedMiddleware is a middleware function that checks if the request is authorized to access the endpoint
*/
func authorizedMiddleware(c *gin.Context) {
	// Check if the Authorization header is present and not empty
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		err := Errors.Unauthorized()
		Logger.Warning(err.Error())
		c.AbortWithStatusJSON(err.StatusCode, err.Err)
		return
	}

	// Check if the token is valid
	isValidated, validatedToken := Auth.ValidateBearerToken(authHeader)
	if !isValidated {
		err := Errors.ForbiddenInvalidToken()
		Logger.Warning(err.Error())
		c.AbortWithStatusJSON(err.StatusCode, err.Err)
		return
	}

	// Add token retried scopes to the context for later use in the endpoint handler
	c.Set("scopes", validatedToken.Scope)
	c.Set("authenticatedUserId", validatedToken.AuthorizedUserID)

	c.Next()
}

/*
requiredScopes checks if the user has the required scopes to access the endpoint, returning an error if the user is missing any of the required scopes
*/
func requiredScopes(c *gin.Context, requiredScopes []Auth.Scope) error {
	scopesInterface, exists := c.Get("scopes")
	if !exists {
		return fmt.Errorf("scopes not found in context")
	}

	scopesStr, ok := scopesInterface.(string)
	if !ok {
		return fmt.Errorf("scopes in context is not a string")
	}

	scopeIDs := strings.Split(scopesStr, " ")
	var scopes []Auth.Scope

	for _, scopeID := range scopeIDs {
		scopes = append(scopes, Auth.GetScopeFromID(scopeID))
	}

	// Check if the user has the required scopes
	missingScopes := []Auth.Scope{}
	for _, requiredScope := range requiredScopes {
		found := false
		for _, scope := range scopes {
			if scope.Name == requiredScope.Name {
				found = true
				break
			}
		}
		if !found {
			missingScopes = append(missingScopes, requiredScope)
		}
	}

	if len(missingScopes) > 0 {
		err := Errors.ForbiddenNoRequiredScope(missingScopes)
		Logger.Warning(err.Error())
		return err
	}

	return nil
}

/*
requiredParamValidation checks if the given parameters (url, query or post body) are present in the request
*/
func requiredParamValidation(c *gin.Context, urlParamsNames []string, queryParamsName []string, postParamsName []string) bool {
	// TODO: to implement
	return true
}

/*
convertInputParamType converts the input parameter from string to the desired type (int, float, bool, etc.)
*/
func convertInputParamType[T any](input string, paramName string) (T, error) {
	// TODO: to implement
	var zeroValue T
	return zeroValue, nil
}

/*
throwError return prettier JSON errors
*/
func throwError(c *gin.Context, err error) {
	if apiError, ok := err.(*Errors.ApiError); ok {
		c.JSON(apiError.StatusCode, apiError.Err)
	} else { // Default error handling
		c.JSON(500, gin.H{"error": "Internal Server Error"})
	}
}

// </editor-fold>
