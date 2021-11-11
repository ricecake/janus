package util

import (
	"bytes"
	"html/template"
	"path/filepath"

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

	message := mail.NewSingleEmail(from, string(subject), to, string(plainTextContent), string(htmlContent))

	client := sendgrid.NewSendClient(key)

	log.Printf("SENDING %+v", message)

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

func RenderEmailTemplate(templateName string, context map[string]interface{}) (subjectOutput, plainOutput, htmlOutput []byte, renderErr error) {
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

	var tpl bytes.Buffer

	if renderErr = emailTemplates.ExecuteTemplate(&tpl, "subject", newContext); renderErr != nil {
		return
	}
	subjectOutput = tpl.Bytes()
	tpl.Reset()

	if renderErr = emailTemplates.ExecuteTemplate(&tpl, "plain", newContext); renderErr != nil {
		return
	}
	plainOutput = tpl.Bytes()
	tpl.Reset()

	if renderErr = emailTemplates.ExecuteTemplate(&tpl, "html", newContext); renderErr != nil {
		return
	}
	htmlOutput = tpl.Bytes()

	return
}
