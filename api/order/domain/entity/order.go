package orderDomainEntity

import (
	"time"

	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
	errorHelper "github.com/ahsansandiah/dpo-test/helpers/error"
	paginateHelper "github.com/ahsansandiah/dpo-test/helpers/paginate"
)

type Order struct {
	ID          int64     `json:"id"`
	CustomerID  int64     `json:"customer_id"`
	OrderDate   time.Time `json:"order_date"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrderRequest struct {
	CustomerID  int64              `json:"customer_id"`
	OrderDate   time.Time          `json:"order_date"`
	TotalAmount float64            `json:"total_amount"`
	OrderItems  []OrderItemRequest `json:"order_items"`
}

type OrderUpdateRequest struct {
	OrderDate   time.Time `json:"order_date"`
	TotalAmount float64   `json:"total_amount"`
}

type OrderResponse struct {
	ID          int64                          `json:"id"`
	OrderDate   time.Time                      `json:"order_date"`
	TotalAmount float64                        `json:"total_amount"`
	Status      string                         `json:"status"`
	Customer    *customerDomainEntity.Customer `json:"customer"`
	Items       []OrderItem                    `json:"items"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
}

type OrderListRespone struct {
	Order  []OrderResponse
	Cursor *paginateHelper.Cursor `json:"Cursor"`
}

type OrderFilter struct {
	CustomerID string `json:"customer_id"`
	OrderDate  string `json:"order_date"`
	Status     string `json:"status"`
	LIMIT      int    `json:"limit"`
	Cursor     string `json:"cursor"`
}

func (r *OrderRequest) Validate() error {
	if r.CustomerID < 0 {
		return errorHelper.ErrorFullNameIsRequired
	}

	if r.OrderDate.IsZero() {
		return errorHelper.ErrorOrderDateRequired
	}

	if r.TotalAmount < 0 {
		return errorHelper.ErrorAmountIsRequired
	}

	if len(r.OrderItems) == 0 {
		return errorHelper.ErrorOrderItemsIsRequired
	}

	return nil
}
