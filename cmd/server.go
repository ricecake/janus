/*
Copyright © 2019 Sebastian Green-Husted <geoffcake@gmail.com>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ricecake/karma_chameleon/http_middleware"
	"github.com/ricecake/janus/model"
	"github.com/ricecake/janus/public_routes"
	"github.com/ricecake/janus/user_routes"
	"github.com/ricecake/karma_chameleon/util"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		ticker := time.NewTicker(10 * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					model.Cleanup()
				case <-quit:
					return
				}
			}
		}()

		util.LoadKey(viper.GetString("security.key"))

		gin.SetMode(viper.GetString("http.mode"))
		ginInterface := viper.GetString("http.interface")
		ginPort := viper.GetInt("http.port")
		ginRunOn := fmt.Sprintf("%s:%d", ginInterface, ginPort)

		r := gin.New()
		setupRouter(r)

		r.Run(ginRunOn)

		close(quit)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func setupRouter(r *gin.Engine) {

	r.Use(http_middleware.RateLimiter())
	r.Use(gin.Logger())
	r.Use(gin.RecoveryWithWriter(log.StandardLogger().Writer()))
	r.Use(http_middleware.SecurityMiddleware())

	rootGroup := r.Group("/")

	staticDir := viper.GetString("http.static")
	if staticDir != "" {
		rootGroup.Static("/static", staticDir)
	}

	contentDir := viper.GetString("http.content")
	if contentDir != "" {
		rootGroup.Static("/content", contentDir)
	}

	// Public routes are special, and need to live outside a namespace
	public_routes.Configure(rootGroup)

	user_routes.Configure(rootGroup.Group("/profile"))
}
