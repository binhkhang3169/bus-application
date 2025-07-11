// ... file: ticket-service/pkg/emailclient/client.go
package emailclient

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"ticket-service/config"
	"ticket-service/domain/models"
	"ticket-service/pkg/kafkaclient"
	"ticket-service/pkg/utils"
)

type EmailClient struct {
	publisher *kafkaclient.Publisher
	logger    utils.Logger
	cfg       config.Config
}

func NewEmailClient(publisher *kafkaclient.Publisher, logger utils.Logger, cfg config.Config) *EmailClient {
	return &EmailClient{
		publisher: publisher,
		logger:    logger,
		cfg:       cfg,
	}
}

func (e *EmailClient) SendStatusEmail(ctx context.Context, email, ticketName sql.NullString, ticketID string, price float64, statusCode string, qrCodeURLs []string) error {
	if !email.Valid || email.String == "" {
		e.logger.Info("Invalid or empty email for ticket %s, cannot send status email.", ticketID)
		return fmt.Errorf("email không hợp lệ hoặc trống cho mã vé: %s", ticketID)
	}

	var title, bodyJSON, emailType string
	var err error

	customerName := "Quý khách"
	if ticketName.Valid && ticketName.String != "" {
		customerName = ticketName.String
	}

	switch statusCode {
	case "1": // Payment success / Ticket active
		title = fmt.Sprintf("Xác nhận vé điện tử thành công - Mã vé %s", ticketID)
		emailType = "ticket_confirmation_payment" // Type mới cho email service

		emailData := models.TicketConfirmationData{
			CustomerName: customerName,
			TicketID:     ticketID,
			Price:        price,
			QRCodeURLs:   qrCodeURLs,
		}
		bodyBytes, marshalErr := json.Marshal(emailData)
		if marshalErr != nil {
			err = fmt.Errorf("failed to marshal ticket confirmation data: %w", marshalErr)
		} else {
			bodyJSON = string(bodyBytes)
		}

	case "2": // Refund success
		title = fmt.Sprintf("Xác nhận hoàn tiền thành công - Mã vé %s", ticketID)
		emailType = "ticket_refund_success" // Type mới

		emailData := models.RefundNotificationData{
			CustomerName: customerName,
			TicketID:     ticketID,
		}
		bodyBytes, marshalErr := json.Marshal(emailData)
		if marshalErr != nil {
			err = fmt.Errorf("failed to marshal refund data: %w", marshalErr)
		} else {
			bodyJSON = string(bodyBytes)
		}

	default:
		e.logger.Info("sendStatusEmail called with unhandled statusCode %s for ticket %s", statusCode, ticketID)
		return fmt.Errorf("status code không hợp lệ (%s) cho việc gửi email thông báo vé %s", statusCode, ticketID)
	}

	eventPayload := kafkaclient.EmailRequestEvent{
		To:    email.String,
		Title: title,
		Body:  bodyJSON,
		Type:  emailType,
	}
	err = e.publisher.Publish(ctx, e.cfg.Kafka.Topics.EmailRequests.Topic, []byte(email.String), eventPayload)
	if err != nil {
		e.logger.Error("Failed to queue email request via Kafka for ticket %s: %v", ticketID, err)
		return err
	}
	return nil
}
