// apps/backend/internal/application/commerce/create_order.go

package commerce

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/commerce"
)

type CreateOrderRequest struct {
	UserID       string
	Items        []OrderItemRequest
	ShippingInfo commerce.ShippingInfo
	Notes        string
}

type OrderItemRequest struct {
	ProductID string
	Quantity  int
}

type CreateOrderUseCase struct {
	orderRepo   commerce.OrderRepository
	productRepo commerce.ProductRepository
}

func NewCreateOrderUseCase(orderRepo commerce.OrderRepository, productRepo commerce.ProductRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{orderRepo: orderRepo, productRepo: productRepo}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, req CreateOrderRequest) (*commerce.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// 1. دریافت اطلاعات محصولات از دیتابیس برای محاسبه قیمت
	var productIDs []string
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
	}
	
	products, err := uc.productRepo.FindProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	
	productMap := make(map[string]*commerce.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	// 2. ساخت اقلام سفارش و محاسبه قیمت کل
	var orderItems []commerce.OrderItem
	var totalCents int
	for _, itemReq := range req.Items {
		product, ok := productMap[itemReq.ProductID]
		if !ok {
			return nil, errors.New("product not found: " + itemReq.ProductID)
		}
		
		totalCents += product.PriceCents * itemReq.Quantity
		orderItems = append(orderItems, commerce.OrderItem{
			ProductID:      product.ID,
			Quantity:       itemReq.Quantity,
			UnitPriceCents: product.PriceCents,
		})
	}
	
	// 3. ایجاد موجودیت سفارش
	newOrder := &commerce.Order{
		ID:            ulid.New().String(),
		UserID:        req.UserID,
		Status:        commerce.StatusPending,
		TotalCents:    totalCents,
		Shipping:      req.ShippingInfo,
		Notes:         req.Notes,
		Items:         orderItems,
		CreatedAt:     time.Now(),
	}
	
	// 4. ذخیره سفارش در دیتابیس
	if err := uc.orderRepo.Create(ctx, newOrder); err != nil {
		return nil, err
	}
	
	return newOrder, nil
}