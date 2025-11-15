// apps/backend/internal/infrastructure/database/postgres/order_repository.go

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"narrative-architecture/apps/backend/internal/domain/commerce"
	"github.com/oklog/ulid/v2"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create یک سفارش جدید را در یک تراکنش ایجاد می‌کند.
func (r *OrderRepository) Create(ctx context.Context, order *commerce.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orderQuery := `
        INSERT INTO orders (id, user_id, status, total_cents, shipping_cents, discount_cents, currency,
                          shipping_name, shipping_phone, shipping_address_line, city, province, postal_code, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `
	_, err = tx.ExecContext(ctx, orderQuery,
		order.ID, order.UserID, order.Status, order.TotalCents, order.ShippingCents, order.DiscountCents, order.Currency,
		order.Shipping.Name, order.Shipping.Phone, order.Shipping.Address, order.Shipping.City, order.Shipping.Province, order.Shipping.PostalCode, order.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	itemQuery := `
        INSERT INTO order_items (id, order_id, product_id, quantity, unit_price_cents, meta)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	stmt, err := tx.PrepareContext(ctx, itemQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare order item statement: %w", err)
	}
	defer stmt.Close()

	for _, item := range order.Items {
		metaJSON, _ := json.Marshal(item.Meta)
		_, err := stmt.ExecContext(ctx, ulid.New().String(), order.ID, item.ProductID, item.Quantity, item.UnitPriceCents, metaJSON)
		if err != nil {
			return fmt.Errorf("failed to insert order item %s: %w", item.ProductID, err)
		}
	}

	return tx.Commit()
}

// FindByID یک سفارش را به همراه اقلام آن پیدا می‌کند.
func (r *OrderRepository) FindByID(ctx context.Context, id string) (*commerce.Order, error) {
	var order commerce.Order
	var paidAt, canceledAt sql.NullTime
	var notes, shippingName, shippingPhone, shippingAddress, city, province, postalCode sql.NullString
	
	query := `
        SELECT id, user_id, status, total_cents, shipping_cents, discount_cents, currency,
               shipping_name, shipping_phone, shipping_address_line, city, province, postal_code, notes,
               created_at, paid_at, canceled_at
        FROM orders
        WHERE id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.UserID, &order.Status, &order.TotalCents, &order.ShippingCents, &order.DiscountCents, &order.Currency,
		&shippingName, &shippingPhone, &shippingAddress, &city, &province, &postalCode, &notes,
		&order.CreatedAt, &paidAt, &canceledAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Handling nullable fields
	if paidAt.Valid { order.PaidAt = &paidAt.Time }
	if canceledAt.Valid { order.CanceledAt = &canceledAt.Time }
	order.Notes = notes.String
	order.Shipping = commerce.ShippingInfo{
		Name:      shippingName.String,
		Phone:     shippingPhone.String,
		Address:   shippingAddress.String,
		City:      city.String,
		Province:  province.String,
		PostalCode: postalCode.String,
	}

	itemsQuery := `SELECT id, product_id, quantity, unit_price_cents, meta FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item commerce.OrderItem
		var metaJSON []byte
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.UnitPriceCents, &metaJSON); err != nil {
			return nil, err
		}
		json.Unmarshal(metaJSON, &item.Meta)
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

// UpdateStatus وضعیت یک سفارش را آپدیت می‌کند.
func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status commerce.OrderStatus, paidAt *time.Time) error {
	query := `UPDATE orders SET status = $1, paid_at = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, paidAt, id)
	return err
}