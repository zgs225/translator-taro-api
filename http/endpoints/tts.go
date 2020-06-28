package endpoints

import (
	"net/http"
	"time"
	"translator-api/app"
	"translator-api/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	ttsService services.TTSService
	qbox       services.QBox
)

func getQBoxService(app *app.Application) services.QBox {
	if qbox == nil {
		var err error
		qbox, err = services.NewSimpleQBox(
			viper.GetString("qiniu.access_key"),
			viper.GetString("qiniu.secret_key"),
			viper.GetString("qiniu.bucket"),
			app.Logger,
		)
		if err != nil {
			panic(err)
		}
	}
	return qbox
}

func getTTSService(app *app.Application) services.TTSService {
	if ttsService == nil {
		var err error
		set := services.NewRedisSet("tts:polly:files", services.MustGetRedisClient())
		qbox := getQBoxService(app)
		ttsService, err = services.NewAWSPollyTTSService(viper.GetString("polly.region"))
		if err != nil {
			panic(err)
		}
		ttsService = services.NewUploadToQBoxTTSServiceMiddleware(qbox, set)(ttsService)
	}

	return ttsService
}

// NewTTSEndpoint 生成 TTS 接口
func NewTTSEndpoint(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		svc := getTTSService(app)

		text := c.PostForm("text")
		lang := c.DefaultPostForm("lang", "en")

		key, err := svc.Speak(text, lang)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		qb := getQBoxService(app)
		url := qb.MakePrivateURL(key, viper.GetString("qiniu.cdn_domain"), time.Hour)

		c.JSON(http.StatusOK, gin.H{
			"code":      0,
			"message":   "OK",
			"timestamp": time.Now().Unix(),
			"data": gin.H{
				"speak_url": url,
			},
		})
	}
}
