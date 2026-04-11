package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "quay-go-api/docs"
)

func healthController() {
	engine.GET("/health", getHealth)
}

// getHealth Get the API health status
// @Summary Get the API health status
// @Description Get the API health status
// @Tags Health
// @Success 200 {string} string "OK"
// @Router /health [get]
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
