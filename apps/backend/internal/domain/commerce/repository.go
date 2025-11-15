// apps/backend/internal/domain/commerce/repository.go

package commerce

import "context"

// OrderRepository رابطی برای عملیات مربوط به سفارش‌ها است.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus, paidAt *time.Time) error
}

// ProductRepository رابطی برای عملیات مربوط به محصولات است.
type ProductRepository interface {
	FindProductsByIDs(ctx context.Context, ids []string) ([]*Product, error)
}

// PaymentRepository رابطی برای عملیات مربوط به پرداخت‌ها است.
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	FindByAuthority(ctx context.Context, authority string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
}