package orderRepository

import (
	"context"
	"database/sql"
	"time"

	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
	orderDomainInterface "github.com/ahsansandiah/dpo-test/api/order/domain"
	orderDomainEntity "github.com/ahsansandiah/dpo-test/api/order/domain/entity"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
)

type Order struct {
	DB  *sql.DB
	log log.Log
	cfg *config.Config
}

func NewOrderRepository(mgr manager.Manager) orderDomainInterface.OrderRepository {
	repo := new(Order)
	repo.DB = mgr.GetDB()
	repo.log = mgr.GetLog()
	repo.cfg = mgr.GetConfig()

	return repo
}

func (r *Order) GetAll(ctx context.Context, filter *orderDomainEntity.OrderFilter) ([]orderDomainEntity.OrderResponse, error) {
	query := `SELECT 
                o.id, o.customer_id, o.order_date, o.status, o.total_amount, o.created_at, o.updated_at,
                c.id, c.full_name, c.address, c.phone_number, c.email, c.is_active, c.created_at, c.updated_at,
                oi.id, oi.order_id, oi.product_name, oi.quantity, oi.price, oi.total_price, oi.created_at, oi.updated_at
              FROM orders o
              INNER JOIN customers c ON o.customer_id = c.id
              LEFT JOIN order_items oi ON o.id = oi.order_id
              WHERE o.deleted_at IS NULL`

	var args []interface{}

	// Apply filters
	if filter != nil {
		if filter.CustomerID != "" {
			query += " AND o.customer_id = ?"
			args = append(args, filter.CustomerID)
		}
		if filter.OrderDate != "" {
			query += " AND o.order_date = ?"
			args = append(args, filter.OrderDate)
		}
		if filter.Status != "" {
			query += " AND o.status = ?"
			args = append(args, filter.Status)
		}
	}

	// Add pagination
	query += " LIMIT ?"
	args = append(args, filter.LIMIT)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[int64]*orderDomainEntity.OrderResponse)
	for rows.Next() {
		var order orderDomainEntity.OrderResponse
		var customer customerDomainEntity.Customer
		var item orderDomainEntity.OrderItem

		// Initialize Customer pointer
		order.Customer = &customer

		err := rows.Scan(
			&order.ID, &order.Customer.ID, &order.OrderDate, &order.Status, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt,
			&order.Customer.ID, &order.Customer.FullName, &order.Customer.Address, &order.Customer.PhoneNumber, &order.Customer.Email, &order.Customer.IsActive, &order.Customer.CreatedAt, &order.Customer.UpdatedAt,
			&item.ID, &item.OrderID, &item.ProductName, &item.Quantity, &item.Price, &item.TotalPrice, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			r.log.ErrorLog(ctx, err)
			return nil, err
		}

		if existingOrder, exists := orderMap[order.ID]; exists {
			existingOrder.Items = append(existingOrder.Items, item)
		} else {
			if item.ID != 0 {
				order.Items = append(order.Items, item)
			}
			orderMap[order.ID] = &order
		}
	}
	if err := rows.Err(); err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	orders := make([]orderDomainEntity.OrderResponse, 0, len(orderMap))
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (r *Order) GetById(ctx context.Context, ID int64) (*orderDomainEntity.Order, error) {
	order := orderDomainEntity.Order{}

	query := "SELECT id, customer_id, order_date, status, total_amount, created_at, updated_at FROM orders WHERE id = ?"
	err := r.DB.QueryRow(query, ID).Scan(&order.ID, &order.CustomerID, &order.OrderDate, &order.Status, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &order, nil
}

func (r *Order) Delete(ctx context.Context, ID int64) error {
	stmt, err := r.DB.Prepare("UPDATE orders SET deleted_at = ? WHERE id = ?")
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return err
	}
	defer stmt.Close()

	// Set deleted_at to current timestamp
	_, err = stmt.Exec(time.Now(), ID)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return err
	}

	return nil
}

func (r *Order) Update(ctx context.Context, ID int64, request *orderDomainEntity.OrderUpdateRequest) (*orderDomainEntity.Order, error) {
	var order orderDomainEntity.Order
	stmt, err := r.DB.PrepareContext(ctx, "UPDATE orders SET order_date = ?, total_amount = ? WHERE id = ?")
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}
	defer stmt.Close()

	// Execute UPDATE statement
	_, err = stmt.ExecContext(ctx, request.OrderDate, request.TotalAmount, ID)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	// Query the updated order
	query := "SELECT id, customer_id, order_date, status, total_amount, created_at, updated_at FROM orders WHERE id = ?"
	err = r.DB.QueryRowContext(ctx, query, ID).Scan(&order.ID, &order.CustomerID, &order.OrderDate, &order.Status, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &order, nil
}

func (r *Order) Create(ctx context.Context, request *orderDomainEntity.OrderRequest) error {
	// Start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	result, err := r.DB.ExecContext(ctx, "INSERT INTO orders (customer_id, order_date, total_amount) VALUES (?, ?, ?)", request.CustomerID, request.OrderDate, request.TotalAmount)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return err
	}

	// Get ID of the inserted record
	orderID, err := result.LastInsertId()
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return err
	}

	// Insert order items
	for _, item := range request.OrderItems {
		_, err := tx.Exec("INSERT INTO order_items (order_id, product_name, quantity, price, total_price) VALUES (?, ?, ?, ?, ?)",
			orderID, item.ProductName, item.Quantity, item.Price, item.TotalPrice)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *Order) GetCustomer(ctx context.Context, customerID int64) (*customerDomainEntity.Customer, error) {
	customer := customerDomainEntity.Customer{}

	query := "SELECT id, full_name, address, phone_number, email, is_active, created_at, updated_at FROM customers WHERE id = ?"
	err := r.DB.QueryRow(query, customerID).Scan(&customer.ID, &customer.FullName, &customer.Address, &customer.PhoneNumber, &customer.Email, &customer.IsActive, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &customer, nil
}

func (r *Order) GetOrderItems(ctx context.Context, orderId int64) ([]orderDomainEntity.OrderItem, error) {
	var args []interface{}

	query := "SELECT id, product_name, price, quantity, total_price FROM order_items"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}
	defer rows.Close()

	var orderItems []orderDomainEntity.OrderItem
	for rows.Next() {
		var item orderDomainEntity.OrderItem
		if err := rows.Scan(&item.ID, &item.ProductName, &item.Price, &item.TotalPrice, &item.Quantity); err != nil {
			r.log.ErrorLog(ctx, err)
			return nil, err
		}
		orderItems = append(orderItems, item)
	}
	if err = rows.Err(); err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return orderItems, nil
}
