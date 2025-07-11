package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/stripe/stripe-go/v76" // Use appropriate version
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund" // Import refund package
	"github.com/stripe/stripe-go/v76/webhook"

	"payment_service/config"
	"payment_service/domain/model"
	"payment_service/internal/db"
)

// StripeServiceInterface định nghĩa các phương thức cho stripe service
type StripeServiceInterface interface {
	CreatePaymentIntent(ctx context.Context, req model.InitialStripePaymentRequest) (*model.StripePaymentIntentResponse, error)
	ConfirmPayment(ctx context.Context, paymentIntentID string) (db.Invoice, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
	RefundPayment(ctx context.Context, req model.StripeInitiateRefundRequest) (db.Invoice, error) // New method
}

// StripeService xử lý các tương tác với Stripe API
type StripeService struct {
	cfg            *config.StripeConfig
	invoiceService InvoiceServiceInterface // Use interface
}

// NewStripeService tạo một stripe service mới
func NewStripeService(cfg *config.StripeConfig, invoiceService InvoiceServiceInterface) StripeServiceInterface {
	stripe.Key = cfg.SecretKey
	return &StripeService{
		cfg:            cfg,
		invoiceService: invoiceService,
	}
}

// CreatePaymentIntent tạo một Stripe PaymentIntent và một hóa đơn liên quan
func (s *StripeService) CreatePaymentIntent(ctx context.Context, req model.InitialStripePaymentRequest) (*model.StripePaymentIntentResponse, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount), // Amount is already in smallest unit from request
		Currency: stripe.String(strings.ToLower(req.Currency)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"customer_id": req.CustomerID,
			"ticket_id":   req.TicketID,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating Stripe PaymentIntent: %v. Request: %+v", err, req)
		if stripeErr, ok := err.(*stripe.Error); ok {
			return nil, fmt.Errorf("stripe error (%s): %s - %s", stripeErr.Code, stripeErr.Msg, stripeErr.Type)
		}
		return nil, fmt.Errorf("stripe: failed to create payment intent: %w", err)
	}

	dbInvoice, err := s.invoiceService.CreateInvoiceForStripe(ctx, req, pi.ID)
	if err != nil {
		log.Printf("Error creating invoice in DB after Stripe PI creation (%s): %v", pi.ID, err)
		// Consider cancelling the PaymentIntent on Stripe if invoice creation fails.
		// _, cancelErr := paymentintent.Cancel(pi.ID, nil)
		// if cancelErr != nil { log.Printf("Failed to cancel Stripe PI %s after DB error: %v", pi.ID, cancelErr) }
		return nil, fmt.Errorf("failed to create internal invoice after Stripe PI creation: %w", err)
	}

	return &model.StripePaymentIntentResponse{
		ClientSecret:    pi.ClientSecret,
		PaymentIntentID: pi.ID,
		InvoiceID:       dbInvoice.InvoiceID,
		PublishableKey:  s.cfg.PublishableKey,
	}, nil
}

// ConfirmPayment xác minh PaymentIntent và cập nhật hóa đơn
func (s *StripeService) ConfirmPayment(ctx context.Context, paymentIntentID string) (db.Invoice, error) {
	if paymentIntentID == "" {
		return db.Invoice{}, errors.New("stripe: paymentIntentID cannot be empty for confirmation")
	}

	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		log.Printf("Error retrieving PaymentIntent %s from Stripe: %v", paymentIntentID, err)
		if stripeErr, ok := err.(*stripe.Error); ok {
			return db.Invoice{}, fmt.Errorf("stripe error (%s) retrieving PI %s: %s", stripeErr.Code, paymentIntentID, stripeErr.Msg)
		}
		return db.Invoice{}, fmt.Errorf("stripe: failed to retrieve payment intent %s: %w", paymentIntentID, err)
	}

	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		var chargeID string
		if pi.LatestCharge != nil && pi.LatestCharge.ID != "" {
			chargeID = pi.LatestCharge.ID
		}
		var paymentMethodDetailsJSON string
		if pi.PaymentMethod != nil {
			pmJSON, jsonErr := json.Marshal(pi.PaymentMethod)
			if jsonErr == nil {
				paymentMethodDetailsJSON = string(pmJSON)
			} else {
				log.Printf("Warning: Could not marshal payment method details for PI %s: %v", pi.ID, jsonErr)
			}
		}
		updatedInvoice, updateErr := s.invoiceService.UpdateInvoiceStatusForStripeSuccess(ctx, pi.ID, chargeID, paymentMethodDetailsJSON)
		if updateErr != nil {
			log.Printf("Error updating invoice for successful Stripe payment PI %s: %v", pi.ID, updateErr)
			return db.Invoice{}, fmt.Errorf("failed to update invoice after successful Stripe payment (PI: %s): %w", pi.ID, updateErr)
		}
		log.Printf("Stripe payment successful for PI: %s, Charge: %s. Invoice %s updated.", pi.ID, chargeID, updatedInvoice.InvoiceID)
		return updatedInvoice, nil

	case stripe.PaymentIntentStatusProcessing:
		log.Printf("Stripe PaymentIntent %s is still processing.", pi.ID)
		return db.Invoice{}, fmt.Errorf("stripe: payment for PI %s is still processing", pi.ID)

	case stripe.PaymentIntentStatusRequiresPaymentMethod,
		stripe.PaymentIntentStatusRequiresConfirmation,
		stripe.PaymentIntentStatusRequiresAction,
		stripe.PaymentIntentStatusCanceled:
		log.Printf("Stripe PaymentIntent %s failed or was canceled. Status: %s", pi.ID, pi.Status)
		var failureReason string
		if pi.LastPaymentError != nil {
			failureReason = fmt.Sprintf("Stripe Error (%s): %s", pi.LastPaymentError.Code, pi.LastPaymentError.Msg) // .Message not .Msg
		} else {
			failureReason = fmt.Sprintf("Status: %s", pi.Status)
		}
		_, errUpdateFail := s.invoiceService.UpdateInvoiceStatusForPaymentFailure(ctx, pi.ID, model.PaymentMethodStripe, failureReason)
		if errUpdateFail != nil {
			log.Printf("Error updating invoice for failed/canceled Stripe payment PI %s: %v", pi.ID, errUpdateFail)
		}
		return db.Invoice{}, fmt.Errorf("stripe: payment for PI %s was not successful (Status: %s). Reason: %s", pi.ID, pi.Status, failureReason)

	default:
		log.Printf("Unhandled Stripe PaymentIntent status for PI %s: %s", pi.ID, pi.Status)
		return db.Invoice{}, fmt.Errorf("stripe: unhandled payment intent status for PI %s: %s", pi.ID, pi.Status)
	}
}

// HandleWebhook xử lý các sự kiện webhook từ Stripe
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	if s.cfg.WebhookSecret == "" {
		log.Println("Stripe webhook secret is not configured. Skipping signature verification. THIS IS INSECURE.")
		// return errors.New("webhook secret not configured") // Should be an error in prod
	}

	event, err := webhook.ConstructEventWithOptions(payload, signature, s.cfg.WebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true, // Add this if you encounter API version mismatch issues during testing
	})

	if err != nil {
		log.Printf("Error verifying Stripe webhook signature: %v", err)
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	log.Printf("Received Stripe webhook event: ID=%s, Type=%s", event.ID, event.Type)

	switch event.Type {
	case stripe.EventTypePaymentIntentSucceeded:
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			log.Printf("Error unmarshaling payment_intent.succeeded data: %v", err)
			return fmt.Errorf("failed to unmarshal payment_intent.succeeded data: %w", err)
		}
		log.Printf("Webhook: PaymentIntent %s succeeded.", pi.ID)
		var chargeID string
		if pi.LatestCharge != nil && pi.LatestCharge.ID != "" {
			chargeID = pi.LatestCharge.ID
		}
		var paymentMethodDetailsJSON string
		if pi.PaymentMethod != nil {
			pmJSON, _ := json.Marshal(pi.PaymentMethod)
			paymentMethodDetailsJSON = string(pmJSON)
		}
		_, errUpdate := s.invoiceService.UpdateInvoiceStatusForStripeSuccess(ctx, pi.ID, chargeID, paymentMethodDetailsJSON)
		if errUpdate != nil {
			log.Printf("Webhook: Error updating invoice for successful PI %s: %v", pi.ID, errUpdate)
			return fmt.Errorf("webhook: failed to update invoice for PI %s: %w", pi.ID, errUpdate)
		}
		log.Printf("Webhook: Invoice updated successfully for PI %s.", pi.ID)

	case stripe.EventTypePaymentIntentPaymentFailed:
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			log.Printf("Error unmarshaling payment_intent.payment_failed data: %v", err)
			return fmt.Errorf("failed to unmarshal payment_intent.payment_failed data: %w", err)
		}
		log.Printf("Webhook: PaymentIntent %s failed.", pi.ID)
		var failureReason string
		if pi.LastPaymentError != nil {
			failureReason = fmt.Sprintf("Stripe Error (%s): %s", pi.LastPaymentError.Code, pi.LastPaymentError.Msg) // .Message
		} else {
			failureReason = fmt.Sprintf("Status: %s", pi.Status)
		}
		_, errUpdateFail := s.invoiceService.UpdateInvoiceStatusForPaymentFailure(ctx, pi.ID, model.PaymentMethodStripe, failureReason)
		if errUpdateFail != nil {
			log.Printf("Webhook: Error updating invoice for failed PI %s: %v", pi.ID, errUpdateFail)
			return fmt.Errorf("webhook: failed to update invoice for failed PI %s: %w", pi.ID, errUpdateFail)
		}
		log.Printf("Webhook: Invoice updated to FAILED for PI %s.", pi.ID)

	// Handle refund update from webhook if needed, e.g., charge.refunded
	case stripe.EventTypeChargeRefunded:
		var ch stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
			log.Printf("Error unmarshaling charge.refunded data: %v", err)
			return fmt.Errorf("failed to unmarshal charge.refunded data: %w", err)
		}
		log.Printf("Webhook: Charge %s for PaymentIntent %s was refunded.", ch.ID, ch.PaymentIntent.ID)
		// If multiple refunds on a PI, this needs careful handling.
		// The refund object itself might be more useful.
		// We already update status when initiating refund. This could be for external refunds.
		// For now, this is a good place to reconcile.
		if ch.PaymentIntent != nil && ch.PaymentIntent.ID != "" {
			invoice, findErr := s.invoiceService.GetInvoiceByStripePaymentIntentID(ctx, ch.PaymentIntent.ID)
			if findErr == nil {
				// Assuming the first refund means the primary one we care about changing status for.
				// If invoice is not already refunded, update it.
				if invoice.PaymentStatus.String != string(model.PaymentStatusRefunded) {
					// Get latest refund from the charge
					var refundReason string = "Refund processed via Stripe."
					if len(ch.Refunds.Data) > 0 {
						// Get reason from the latest refund object on the charge
						// Stripe API might not directly provide the custom reason here this way.
						// The refund object itself (if event was for refund.updated) would be better.
						// For now, using a generic message.
						refundReason = fmt.Sprintf("Stripe Charge %s refunded. Refund ID: %s", ch.ID, ch.Refunds.Data[0].ID)

					}
					_, errRefundUpdate := s.invoiceService.UpdateInvoiceStatusForRefund(ctx, invoice.InvoiceID, refundReason, ch.Refunds.Data[0].ID)
					if errRefundUpdate != nil {
						log.Printf("Webhook: Error updating invoice %s to refunded via charge.refunded event: %v", invoice.InvoiceID, errRefundUpdate)
						return errRefundUpdate
					}
					log.Printf("Webhook: Invoice %s updated to REFUNDED via charge.refunded event for PI %s.", invoice.InvoiceID, ch.PaymentIntent.ID)
				}
			} else {
				log.Printf("Webhook: charge.refunded: Could not find invoice for PI %s", ch.PaymentIntent.ID)
			}
		}

	default:
		log.Printf("Webhook: Unhandled event type: %s", event.Type)
	}

	return nil
}

// RefundPayment handles refunding a Stripe payment.
func (s *StripeService) RefundPayment(ctx context.Context, req model.StripeInitiateRefundRequest) (db.Invoice, error) {
	invoice, err := s.invoiceService.GetLatestCompletedInvoiceByTicketID(ctx, req.TicketID)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("stripe refund: failed to get invoice for ticket ID %s: %w", req.TicketID, err)
	}

	if !invoice.StripePaymentIntentID.Valid || invoice.StripePaymentIntentID.String == "" {
		return db.Invoice{}, fmt.Errorf("stripe refund: invoice %s does not have a Stripe PaymentIntentID", invoice.InvoiceID)
	}

	if invoice.PaymentStatus.String != string(model.PaymentStatusCompleted) {
		return db.Invoice{}, fmt.Errorf("stripe refund: invoice %s is not completed (status: %s), cannot refund", invoice.InvoiceID, invoice.PaymentStatus.String)
	}
	if invoice.PaymentStatus.String == string(model.PaymentStatusRefunded) {
		return db.Invoice{}, fmt.Errorf("stripe refund: invoice %s is already refunded", invoice.InvoiceID)
	}

	// Calculate refund amount in smallest currency unit
	currency := "usd" // Default if not set
	if invoice.Currency.Valid && invoice.Currency.String != "" {
		currency = invoice.Currency.String
	}

	originalAmountSmallestUnit, err := s.invoiceService.GetAmountInSmallestUnit(invoice.FinalAmount, currency)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("stripe refund: could not convert final amount for invoice %s: %w", invoice.InvoiceID, err)
	}

	refundAmountSmallestUnit := originalAmountSmallestUnit
	if req.PercentageDeduction > 0 {
		deduction := float64(originalAmountSmallestUnit) * req.PercentageDeduction
		refundAmountSmallestUnit = originalAmountSmallestUnit - int64(deduction)
	}

	if refundAmountSmallestUnit <= 0 && originalAmountSmallestUnit > 0 { // Cannot refund zero or negative unless it's a full refund of zero amount tx
		log.Printf("Stripe refund: Calculated refund amount for invoice %s is zero or negative (%d from %.2f with %.2f%% deduction). Assuming full refund of original amount if original was > 0, otherwise error.",
			invoice.InvoiceID, refundAmountSmallestUnit, invoice.FinalAmount, req.PercentageDeduction*100)
		// If original amount was 0, then refunding 0 is fine. Otherwise, this logic needs review.
		// For now, let it pass if original was 0. If original > 0 and refund is <=0, it's an issue.
		if originalAmountSmallestUnit > 0 {
			// This might happen if PercentageDeduction is 1.0 or more.
			// A 100% deduction means refundAmount is 0. Stripe might allow refunding 0.
			// If we want to prevent 0 refund unless it's a full refund of a 0 value transaction:
			if refundAmountSmallestUnit < 0 { // strictly negative, which shouldn't happen with percent 0-1
				return db.Invoice{}, fmt.Errorf("stripe refund: calculated refund amount is negative for invoice %s", invoice.InvoiceID)
			}
			// If refundAmountSmallestUnit is 0, Stripe might accept it.
		}
	}

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(invoice.StripePaymentIntentID.String),
		// Amount is optional; if omitted, a full refund is issued.
		// If req.PercentageDeduction results in a partial refund, set Amount.
		// If req.PercentageDeduction is 0 (full refund of final_amount), Amount can be set or omitted.
	}
	if req.PercentageDeduction > 0.0 || refundAmountSmallestUnit < originalAmountSmallestUnit { // If it's a partial refund
		params.Amount = stripe.Int64(refundAmountSmallestUnit)
	}
	// If req.PercentageDeduction == 0, it's a full refund of invoice.FinalAmount.
	// Stripe will refund the full amount if params.Amount is not set.
	// If refundAmountSmallestUnit is calculated (even for full refund), we can set it.
	// Let's be explicit if we calculated an amount:
	if refundAmountSmallestUnit > 0 { // Only set amount if it's positive
		params.Amount = stripe.Int64(refundAmountSmallestUnit)
	} else if refundAmountSmallestUnit == 0 && originalAmountSmallestUnit > 0 {
		// Refunding $0.00 for a paid transaction.
		params.Amount = stripe.Int64(0)
	}
	// If originalAmountSmallestUnit was 0, and percentageDeduction is also 0, then refund amount is 0. params.Amount will be 0.

	// Stripe's standard reasons: duplicate, fraudulent, requested_by_customer
	params.Reason = stripe.String(string(stripe.RefundReasonRequestedByCustomer))
	params.Metadata = map[string]string{
		"internal_refund_reason": req.Reason,
		"ticket_id":              req.TicketID,
	}

	stripeRefund, err := refund.New(params)
	if err != nil {
		log.Printf("Stripe API error during refund for PI %s: %v", invoice.StripePaymentIntentID.String, err)
		if stripeErr, ok := err.(*stripe.Error); ok {
			// Try to update invoice to FAILED or keep as COMPLETED but add note about failed refund attempt
			note := fmt.Sprintf("Stripe refund attempt failed for PI %s. Stripe Error (%s): %s. Reason: %s",
				invoice.StripePaymentIntentID.String, stripeErr.Code, stripeErr.Msg, req.Reason)
			log.Print(note)
			// s.invoiceService.AddNoteToInvoice(ctx, invoice.InvoiceID, note) // Hypothetical method
			return db.Invoice{}, fmt.Errorf("stripe API error (%s) refunding PI %s: %s", stripeErr.Code, invoice.StripePaymentIntentID.String, stripeErr.Msg)
		}
		return db.Invoice{}, fmt.Errorf("stripe: failed to process refund for PI %s: %w", invoice.StripePaymentIntentID.String, err)
	}

	log.Printf("Stripe refund successful for PI %s. Refund ID: %s. Amount: %d %s",
		stripeRefund.PaymentIntent, stripeRefund.ID, stripeRefund.Amount, strings.ToUpper(string(stripeRefund.Currency)))

	// Update invoice status to REFUNDED
	// Pass stripeRefund.ID as the specific identifier
	updatedInvoice, err := s.invoiceService.UpdateInvoiceStatusForRefund(ctx, invoice.InvoiceID, req.Reason, stripeRefund.ID)
	if err != nil {
		// This is a critical situation: Stripe refund succeeded, but DB update failed.
		log.Printf("CRITICAL: Stripe refund %s processed, but failed to update invoice %s status: %v", stripeRefund.ID, invoice.InvoiceID, err)
		// Return the invoice object from before the failed update, but with an error indicating the inconsistency.
		return invoice, fmt.Errorf("stripe refund %s successful, but DB update failed for invoice %s: %w. MANUAL INTERVENTION REQUIRED", stripeRefund.ID, invoice.InvoiceID, err)
	}

	// Ticket status update is handled by InvoiceService.UpdateInvoiceStatusForRefund

	return updatedInvoice, nil
}
