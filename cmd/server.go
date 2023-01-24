/*
Copyright © 2019 Sebastian Green-Husted <geoffcake@gmail.com>
*/
package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"janus/admin_routes"
	"janus/model"
	"janus/public_routes"
	"janus/user_routes"
	"janus/util"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/karma_chameleon/http_middleware"
	kutil "github.com/ricecake/karma_chameleon/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		kutil.LoadKey(viper.GetString("security.key"))

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

	tmpl := template.New("react")
	tmpl.Funcs(template.FuncMap{
		"json": func(data gin.H) template.JS {
			output, err := json.Marshal(data)
			if err != nil {
				log.Error(err)
			}
			return template.JS(string(output))
		},
	})

	r.SetHTMLTemplate(tmpl)

	walkErr := fs.WalkDir(util.Content, "content", func(currPath string, dir fs.DirEntry, err error) error {
		if path.Ext(currPath) == ".html" {
			relative, relErr := filepath.Rel("content", currPath)
			if relErr != nil {
				return relErr
			}
			tplData, readErr := fs.ReadFile(util.Content, currPath)
			if readErr != nil {
				return readErr
			}

			_, parseErr := tmpl.New(relative).Parse(string(tplData))
			return parseErr
		}
		return nil
	})

	if walkErr != nil {
		log.Fatal(walkErr)
	}

	staticFs := EmbedFolder(util.Content, "content")
	r.Use(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		if found, index := staticFs.Exists("/", c.Request.URL.Path); found {
			if index {
				reqPath = path.Join(reqPath, "index.html")
			}

			if path.Ext(reqPath) == ".html" {
				c.HTML(200, strings.TrimPrefix(reqPath, "/"), gin.H{
					"CspNonce": c.GetString("CspNonce"),
				})
			} else {
				c.FileFromFS(reqPath, staticFs)
			}
			c.Abort()
		}
	})

	rootGroup := r.Group("/")

	// Public routes are special, and need to live outside a namespace
	public_routes.Configure(rootGroup)

	user_routes.Configure(rootGroup.Group("/profile"))
	admin_routes.Configure(rootGroup.Group("/admin"))
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, reqPath string) (found, index bool) {
	if reqPath != "/" {
		reqPath = strings.TrimSuffix(reqPath, "/")
	}

	file, err := e.Open(reqPath)
	if err != nil {
		return false, false
	}

	stats, err := file.Stat()
	if err != nil {
		return false, false
	}

	isIndex := false
	if stats.IsDir() {
		index := path.Join(reqPath, "index.html")
		_, err := e.Open(index)
		if err != nil {
			return false, false
		}
		isIndex = true
	}

	return true, isIndex
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) embedFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
