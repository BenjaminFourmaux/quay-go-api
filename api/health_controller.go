package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func healthController() {
	health := engine.Group("/health")
	{
		health.GET("", getHealth)
	}
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
