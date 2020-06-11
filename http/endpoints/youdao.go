package endpoints

import (
	"fmt"
	"math"
	"net/http"
	"translator-api/app"
	"translator-api/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang/groupcache/lru"
	"github.com/spf13/viper"
)

var _youdaoSvc services.YoudaoService

func getYoudaoService(app *app.Application) services.YoudaoService {
	if _youdaoSvc == nil {
		appID := viper.GetString("youdao.app_id")
		appKey := viper.GetString("youdao.app_key")
		_youdaoSvc = services.NewBasicYoudaoService(appID, appKey)

		if viper.GetBool("cache.redis") {
			redisClient := redis.NewClient(&redis.Options{
				Addr:         fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
				DB:           viper.GetInt("redis.db"),
				Password:     viper.GetString("redis.password"),
				PoolSize:     50,
				MinIdleConns: 25,
			})

			if err := redisClient.Ping().Err(); err != nil {
				app.Logger.Panic(err)
			}

			_youdaoSvc = services.NewRedisCachedYoudaoServiceMiddleware(redisClient, app)(_youdaoSvc)
		}

		if viper.GetBool("cache.lru") {
			cache := lru.New(int(math.Max(100, viper.GetFloat64("lru.size"))))
			_youdaoSvc = services.NewMemoryCachedYoudaoService(cache, app)(_youdaoSvc)
		}
	}
	return _youdaoSvc
}

// YoudaoEndpoints 有道词典接口
type YoudaoEndpoints struct {
	App *app.Application
}

// CreateQueryEndpoint 生成查询接口
func (edp *YoudaoEndpoints) CreateQueryEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		if len(q) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if len(q) > 255 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		result, err := getYoudaoService(edp.App).Query(q)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
