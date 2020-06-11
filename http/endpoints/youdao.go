package endpoints

import (
	"net/http"
	"translator-api/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var _youdaoSvc services.YoudaoService

func getYoudaoService() services.YoudaoService {
	if _youdaoSvc == nil {
		appID := viper.GetString("youdao.app_id")
		appKey := viper.GetString("youdao.app_key")
		_youdaoSvc = services.NewBasicYoudaoService(appID, appKey)
	}
	return _youdaoSvc
}

// EndpointYoudaoQuery 有道查询接口
func EndpointYoudaoQuery(c *gin.Context) {
	q := c.Query("q")
	if len(q) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(q) > 255 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := getYoudaoService().Query(q)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
