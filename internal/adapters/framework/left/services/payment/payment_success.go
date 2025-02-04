package payment

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strings"
	"text/template"

	"ecom-api/pkg/configs"
	"ecom-api/utils"
)

func sendHtmlEmail(to []string, subject string, htmlBody string) error {
	auth := smtp.PlainAuth(
		"",
		configs.Envs.FromEmail,
		configs.Envs.FromEmailPassword,
		configs.Envs.FromEmailSMTP,
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody
	return smtp.SendMail(
		configs.Envs.SMTPAddress,
		auth,
		configs.Envs.FromEmail,
		to,
		[]byte(message),
	)
}

func HTMLTemplateEmailHandler(w http.ResponseWriter, r *http.Request, addr string, vars map[string]string) {
	basePathForEmailHtml := "./static/"
	emailSubject := "Purchase Successfull!! 🎉"

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	// Convert Param3 (comma-separated string) to a slice of strings
	to := strings.Split(addr, ",")

	// Parse the HTML template
	templatePath := filepath.Join(basePathForEmailHtml, "purchase_success.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to parse template: %v", err))
		return
	}

	// Render the template with the map data
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, vars); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to render template: %v", err))
		return
	}

	err = sendHtmlEmail(to, emailSubject, rendered.String())
	if err != nil {
		log.Println(err.Error())
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]bool{"emailSent": true}, nil)
}
