// apps/backend/internal/application/commerce/verify_purchase.go

package commerce

import (
	"context"
	"errors"
	"time"

	"narrative-architecture/apps/backend/internal/domain/commerce"
)

type PaymentVerifier interface {
	VerifyPayment(ctx context.Context, amount int, authority string) (*commerce.VerificationResponse, error)
}

type VerifyPurchaseRequest struct {
	OrderID   string
	Authority string
	Status    string // "OK" or "NOK" from Zarinpal
}

type VerifyPurchaseUseCase struct {
	orderRepo   commerce.OrderRepository
	paymentRepo commerce.PaymentRepository
	verifier    PaymentVerifier
	// entitlementRepo, ...
}

// ... NewVerifyPurchaseUseCase ...

func (uc *VerifyPurchaseUseCase) Execute(ctx context.Context, req VerifyPurchaseRequest) error {
	payment, err := uc.paymentRepo.FindByAuthority(ctx, req.Authority)
	if err != nil { return err }
	
	if payment.OrderID != req.OrderID { return errors.New("order mismatch") }
	if payment.Status != commerce.PaymentStatusPending { return errors.New("payment already processed") }
	
	if req.Status != "OK" {
		payment.Status = commerce.PaymentStatusFailed
		uc.orderRepo.UpdateStatus(ctx, req.OrderID, commerce.StatusFailed, nil)
		return uc.paymentRepo.Update(ctx, payment)
	}

	// 1. تایید پرداخت با درگاه
	verifyRes, err := uc.verifier.VerifyPayment(ctx, payment.AmountCents, req.Authority)
	if err != nil {
		payment.Status = commerce.PaymentStatusFailed
		uc.orderRepo.UpdateStatus(ctx, req.OrderID, commerce.StatusFailed, nil)
		uc.paymentRepo.Update(ctx, payment)
		return err
	}

	// 2. آپدیت وضعیت پرداخت و سفارش
	payment.Status = commerce.PaymentStatusSuccess
	payment.RefID = verifyRes.RefID
	if err := uc.paymentRepo.Update(ctx, payment); err != nil { return err }

	now := time.Now()
	if err := uc.orderRepo.UpdateStatus(ctx, req.OrderID, commerce.StatusPaid, &now); err != nil { return err }
	
	// 3. اعطای دسترسی‌های دیجیتال (Entitlements)
	// order, _ := uc.orderRepo.FindByID(ctx, req.OrderID)
	// for _, item := range order.Items {
	// 	if item.Type == "course" || item.Type == "premium_content" {
	// 		uc.entitlementRepo.GrantAccess(...)
	// 	}
	// }
	
	return nil
}