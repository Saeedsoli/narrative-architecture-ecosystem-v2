// apps/backend/internal/domain/commerce/order.go

package commerce

import "time"

type OrderStatus string

const (
	StatusPending  OrderStatus = "pending"
	StatusPaid     OrderStatus = "paid"
	StatusFailed   OrderStatus = "failed"
	StatusCanceled OrderStatus = "canceled"
)

type ShippingInfo struct {
	Name      string
	Phone     string
	Address   string
	City      string
	Province  string
	PostalCode string
}

type OrderItem struct {
	ID             string
	ProductID      string
	Quantity       int
	UnitPriceCents int
	Meta           map[string]interface{}
}

type Order struct {
	ID             string
	UserID         string
	Status         OrderStatus
	TotalCents     int
	ShippingCents  int
	DiscountCents  int
	Currency       string
	Shipping       ShippingInfo
	Notes          string
	Items          []OrderItem
	CreatedAt      time.Time
	PaidAt         *time.Time
	CanceledAt     *time.Time
}