// apps/backend/internal/domain/commerce/payment.go

package commerce

import "time"

type PaymentStatus string

const (
	PaymentStatusInit    PaymentStatus = "init"
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID        string
	OrderID   string
	Gateway   string // "zarinpal"
	AmountCents int
	Status    PaymentStatus
	Authority string // شناسه درگاه برای شروع پرداخت
	RefID     string // شناسه مرجع پس از پرداخت موفق
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentRequestResponse ساختار پاسخ از درگاه پرداخت هنگام ایجاد درخواست است.
type PaymentRequestResponse struct {
	PaymentID  string // شناسه داخلی ما
	PaymentURL string // URL برای هدایت کاربر
	Authority  string
}

// VerificationResponse ساختار پاسخ از درگاه پرداخت هنگام تأیید است.
type VerificationResponse struct {
	RefID string
}