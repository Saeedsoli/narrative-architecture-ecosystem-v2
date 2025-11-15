// apps/backend/internal/application/commerce/process_payment.go

package commerce

import (
	"context"
	"fmt"

	"narrative-architecture/apps/backend/internal/domain/commerce"
)

type PaymentGateway interface {
	CreatePaymentRequest(ctx context.Context, amount int, currency, description, callbackURL string) (*commerce.PaymentRequestResponse, error)
}

type ProcessPaymentRequest struct {
	OrderID string
	UserID  string
}

type ProcessPaymentResponse struct {
	PaymentURL string
	PaymentID  string
}

type ProcessPaymentUseCase struct {
	orderRepo   commerce.OrderRepository
	paymentRepo commerce.PaymentRepository
	gateway     PaymentGateway
}

func NewProcessPaymentUseCase(orderRepo commerce.OrderRepository, paymentRepo commerce.PaymentRepository, gateway PaymentGateway) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
		gateway:     gateway,
	}
}

func (uc *ProcessPaymentUseCase) Execute(ctx context.Context, req ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil { return nil, err }
	
	if order.UserID != req.UserID { return nil, errors.New("user does not own this order") }
	if order.Status != commerce.StatusPending { return nil, errors.New("order is not in pending state") }

	// 1. ایجاد درخواست پرداخت در درگاه
	callbackURL := fmt.Sprintf("https://your-frontend.com/payment/callback?orderId=%s", order.ID)
	desc := fmt.Sprintf("پرداخت سفارش شماره %s", order.ID)
	
	paymentReqRes, err := uc.gateway.CreatePaymentRequest(ctx, order.TotalCents, order.Currency, desc, callbackURL)
	if err != nil { return nil, err }

	// 2. ذخیره اطلاعات پرداخت در دیتابیس
	newPayment := &commerce.Payment{
		ID:         paymentReqRes.PaymentID,
		OrderID:    order.ID,
		Gateway:    "zarinpal",
		AmountCents: order.TotalCents,
		Status:     commerce.PaymentStatusPending,
		Authority:  paymentReqRes.Authority,
	}
	
	if err := uc.paymentRepo.Create(ctx, newPayment); err != nil { return nil, err }

	return &ProcessPaymentResponse{
		PaymentURL: paymentReqRes.PaymentURL,
		PaymentID:  newPayment.ID,
	}, nil
}