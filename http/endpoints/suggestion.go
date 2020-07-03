package endpoints

import (
	"net/http"
	"time"
	"translator-api/app"
	"translator-api/services"

	"github.com/gin-gonic/gin"
)

var _suggestionService services.SuggestionService

func getSuggestionService(app *app.Application) services.SuggestionService {
	if _suggestionService == nil {
		_suggestionService = &services.EcDictSuggestionService{
			Dict: app.Dict,
		}
	}
	return _suggestionService
}

// NewSuggestionEndpoint 获取提示端口
func NewSuggestionEndpoint(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")

		if len(q) < 3 {
			c.JSON(http.StatusOK, gin.H{
				"code": 1001,
			})
			return
		}

		v, err := getSuggestionService(app).Suggest(q)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":      0,
			"message":   "ok",
			"timestamp": time.Now().Unix(),
			"data":      v,
		})
	}
}
