package util

import (
	"github.com/flosch/pongo2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	inited  bool
	loader  *pongo2.LocalFilesystemLoader
	pageSet *pongo2.TemplateSet
)

type TemplateContext map[string]interface{}

func RenderHTMLTemplate(template string, context map[string]interface{}) (output []byte, renderErr error) {
	return RenderTemplate("content", template, context)
}

func RenderEmailTemplate(template string, context map[string]interface{}) (subjectOutput, plainOutput, htmlOutput []byte, renderErr error) {
	subjectPath := filepath.Join(template, "subject")
	plainPath := filepath.Join(template, "plain")
	htmlPath := filepath.Join(template, "html")

	subjectOutput, renderErr = RenderTemplate("email", subjectPath, context)
	if renderErr != nil {
		return
	}
	plainOutput, renderErr = RenderTemplate("email", plainPath, context)
	if renderErr != nil {
		return
	}
	htmlOutput, renderErr = RenderTemplate("email", htmlPath, context)
	return
}

func RenderTemplate(style, template string, context map[string]interface{}) (output []byte, renderErr error) {
	templatePath := filepath.Join(viper.GetString("template.path"), style, template)

	ensureTemplates()

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

func ensureTemplates() {
	if !inited {
		loader = pongo2.MustNewLocalFileSystemLoader(viper.GetString("template.path"))
		pageSet = pongo2.NewSet("Janus", loader)

		pageSet.Debug = viper.GetBool("template.debug")

		for index, element := range viper.GetStringMap("template.globals") {
			log.Info(index, element)
			pageSet.Globals[index] = element
		}

		inited = true
	}
}
