package cmd

import (
	"translator-api/app"
	"translator-api/http/endpoints"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动服务器",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.Default()

		router := app.Group("/v1/api")

		router.GET("/ping", endpoints.EndpointPing)

		{
			edps := &endpoints.YoudaoEndpoints{App: app}
			router.GET("/youdao", edps.CreateQueryEndpoint())
		}

		{
			router.POST("/tts", endpoints.NewTTSEndpoint(app))
		}

		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
