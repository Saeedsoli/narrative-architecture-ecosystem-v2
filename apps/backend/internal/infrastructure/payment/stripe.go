// apps/backend/internal/infrastructure/payment/stripe.go

package payment

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/commerce"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

type StripeGateway struct {
	apiKey string
}

func NewStripeGateway(apiKey string) *StripeGateway {
	stripe.Key = apiKey
	return &StripeGateway{apiKey: apiKey}
}

func (s *StripeGateway) CreatePaymentRequest(ctx context.Context, amount int, currency, description, callbackURL string) (*commerce.PaymentRequestResponse, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:     stripe.String(currency),
					UnitAmount:   stripe.Int64(int64(amount)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(description),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(callbackURL + "?success=true&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(callbackURL + "?canceled=true"),
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return &commerce.PaymentRequestResponse{
		PaymentID:  sess.ID,
		PaymentURL: sess.URL,
	}, nil
}

// ... (متد VerifyPayment برای Stripe)