package orderDomainInterface

import (
	"context"
	"net/http"

	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
	orderDomainEntity "github.com/ahsansandiah/dpo-test/api/order/domain/entity"
)

type OrderHandler interface {
	GetAll() http.Handler
	Delete() http.Handler
	GetByID() http.Handler
	Update() http.Handler
	Create() http.Handler
}

type OrderUsecase interface {
	GetAll(ctx context.Context, filter *orderDomainEntity.OrderFilter) ([]orderDomainEntity.OrderResponse, error)
	Delete(ctx context.Context, ID int64) error
	GetByID(ctx context.Context, ID int64) (*orderDomainEntity.OrderResponse, error)
	Update(ctx context.Context, ID int64, request *orderDomainEntity.OrderUpdateRequest) (*orderDomainEntity.OrderResponse, error)
	Create(ctx context.Context, request *orderDomainEntity.OrderRequest) (*orderDomainEntity.OrderRequest, error)
	ValidateCustomer(ctx context.Context, customerID int64) bool
}

type OrderRepository interface {
	GetAll(ctx context.Context, filter *orderDomainEntity.OrderFilter) ([]orderDomainEntity.OrderResponse, error)
	GetById(ctx context.Context, ID int64) (*orderDomainEntity.Order, error)
	Delete(ctx context.Context, ID int64) error
	Update(ctx context.Context, ID int64, request *orderDomainEntity.OrderUpdateRequest) (*orderDomainEntity.Order, error)
	Create(ctx context.Context, request *orderDomainEntity.OrderRequest) error
	GetCustomer(ctx context.Context, customerID int64) (*customerDomainEntity.Customer, error)
	GetOrderItems(ctx context.Context, orderId int64) ([]orderDomainEntity.OrderItem, error)
}
