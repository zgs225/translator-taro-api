package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// EndpointPing ping 响应
func EndpointPing(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
