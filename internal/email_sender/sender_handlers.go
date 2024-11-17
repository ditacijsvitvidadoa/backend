package email_sender

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"path/filepath"
)

type OrderConfirmationData struct {
	OrderNumber string
	FirstName   string
}

func SendOrderConfirmation(to, firstName, orderNumber string) error {
	subject := fmt.Sprintf("Успішно оформлено замовлення №%s", orderNumber)

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
