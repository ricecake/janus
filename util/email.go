package util

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TemplateContext map[string]interface{}

func SendMail(name, address, template string, context TemplateContext) error {

	key := viper.GetString("email.api_key")

	to := mail.NewEmail(name, address)
	from := mail.NewEmail(viper.GetString("email.sender.name"), viper.GetString("email.sender.address"))

	subject, plainTextContent, htmlContent, err := RenderEmailTemplate(template, context)

	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(key)

	response, err := client.Send(message)
	if err != nil {
		log.Error(err)
	} else {
		log.Info(response.StatusCode)
		log.Info(response.Body)
		log.Info(response.Headers)
	}
	return err
}

func RenderEmailTemplate(templateName string, context map[string]interface{}) (subjectOutput, plainOutput, htmlOutput string, renderErr error) {
	subjectPath := filepath.Join("content", "email", templateName, "subject")
	plainPath := filepath.Join("content", "email", templateName, "plain")
	htmlPath := filepath.Join("content", "email", templateName, "html")

	emailTemplates := template.Must(template.ParseFS(Content, subjectPath, plainPath, htmlPath))

	newContext := make(map[string]interface{})
	for index, element := range viper.GetStringMap("template.globals") {
		newContext[index] = element
	}

	for index, element := range context {
		newContext[index] = element
	}

	var tpl strings.Builder

	if renderErr = emailTemplates.ExecuteTemplate(&tpl, "subject", newContext); renderErr != nil {
		return
	}
	subjectOutput = tpl.String()
	tpl.Reset()

	if renderErr = emailTemplates.ExecuteTemplate(&tpl, "plain", newContext); renderErr != nil {
		return
	}
	plainOutput = tpl.String()
	tpl.Reset()

	if renderErr = emailTemplates.ExecuteTemplate(&tpl, "html", newContext); renderErr != nil {
		return
	}
	htmlOutput = tpl.String()

	return
}
