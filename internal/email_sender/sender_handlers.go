package email_sender

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"path/filepath"
)

type OrderConfirmationData struct {
	OrderNumber int
	FirstName   string
}

type SupportQuestionData struct {
	Name        string
	Phone       string
	Email       string
	Title       string
	Description string
}

func SendOrderConfirmation(to, firstName string, orderNumber int) error {
	subject := fmt.Sprintf("Успішно оформлено замовлення №%d", orderNumber)

	templatePath := filepath.Join("/app", "internal", "email_sender", "template_html", "order_confirmation.html")
	absTemplatePath, err := filepath.Abs(templatePath)
	if err != nil {
		return err
	}

	data := OrderConfirmationData{OrderNumber: orderNumber, FirstName: firstName}
	body, err := utils.ParseTemplate(absTemplatePath, data)
	if err != nil {
		return err
	}
	return EmailSender(to, subject, body)
}

func SendSupportQuestion(name, to, phone, title, description string) error {
	subject := fmt.Sprintf("Новий запит до служби підтримки від %s", name)

	templatePath := filepath.Join("/app", "internal", "email_sender", "template_html", "support_question.html")
	absTemplatePath, err := filepath.Abs(templatePath)
	if err != nil {
		return err
	}

	data := SupportQuestionData{
		Name:        name,
		Phone:       phone,
		Email:       to,
		Title:       title,
		Description: description,
	}

	body, err := utils.ParseTemplate(absTemplatePath, data)
	if err != nil {
		return err
	}

	return EmailSender(to, subject, body)
}
