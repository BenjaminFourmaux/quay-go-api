package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Database"
	_ "quay-go-api/docs"
)

func healthController() {
	engine.GET("/health", getHealth)
}

// getHealth Get the API health status
// @Summary Get the API health status
// @Description Get the API health status
// @Tags Health
// @Success 200 {string} healthResponse
// @Success 503 {string} healthResponse
// @Router /health [get]
func getHealth(c *gin.Context) {
	if err := Database.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, "DEGRADED")
	} else {
		c.JSON(http.StatusOK, "OK")
	}
	return
}
