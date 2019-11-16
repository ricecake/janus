/*
Copyright © 2019 Sebastian Green-Husted <geoffcake@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/public_routes"
	"github.com/ricecake/janus/util"
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
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-quit:
					return
				}
			}
		}()

		util.LoadKey(viper.GetString("security.key"))
		util.InitTemplates()

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

	r.Use(gin.Logger())
	r.Use(gin.RecoveryWithWriter(log.StandardLogger().Writer()))

	rootGroup := r.Group("/")
	staticDir := viper.GetString("http.static")
	if staticDir != "" {
		rootGroup.Static("/static", staticDir)
	}

	// Public routes are special, and need to live outside a namespace
	public_routes.Configure(rootGroup)
}
