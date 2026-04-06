package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"quay-go-api/service/logger"
)

var engine *gin.Engine

func StartServer() {
	engine = gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server on port " + port)

	endpointsRegistration()

	err := engine.Run(":" + port)
	if err != nil {
		logger.Error("Failed to start server: " + err.Error())
		return
	}
}

// <editor-fold desc="Private functions">

func endpointsRegistration() {
	healthController()
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
	// TODO: to implement
	c.JSON(500, gin.H{"error": "Internal Server Error"})
}

// </editor-fold>
