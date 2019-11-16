package util

import (
	"github.com/flosch/pongo2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	loader  *pongo2.LocalFilesystemLoader
	pageSet *pongo2.TemplateSet
)

type TemplateContext map[string]interface{}

func RenderTemplate(template string, context map[string]interface{}) (output []byte, renderErr error) {
	templatePath := filepath.Join(viper.GetString("http.template.path"), "content", template)

	templateBody, templateError := pageSet.FromCache(templatePath)
	if templateError != nil {
		renderErr = templateError
		return
	}

	output, renderErr = templateBody.ExecuteBytes(context)

	if renderErr != nil {
		log.Error(renderErr)
	}

	return
}

func InitTemplates() {
	loader = pongo2.MustNewLocalFileSystemLoader(viper.GetString("http.template.path"))
	pageSet = pongo2.NewSet("html", loader)

	pageSet.Debug = viper.GetBool("http.template.debug")
	pageSet.Globals = viper.GetStringMap("http.template.globals")
}
