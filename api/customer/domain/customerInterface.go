package customerDomainInterface

import (
	"context"
	"net/http"

	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
)

type CustomerHandler interface {
	GetAll() http.Handler
	Delete() http.Handler
	GetByID() http.Handler
	Update() http.Handler
	Create() http.Handler
}

type CustomerUsecase interface {
	GetAll(ctx context.Context, filter *customerDomainEntity.CustomerFilter) ([]customerDomainEntity.Customer, error)
	Delete(ctx context.Context, ID int64) error
	GetByID(ctx context.Context, ID int64) (*customerDomainEntity.Customer, error)
	Update(ctx context.Context, ID int64, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error)
	Create(ctx context.Context, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error)
}

type CustomerRepository interface {
	GetAll(ctx context.Context, filter *customerDomainEntity.CustomerFilter) ([]customerDomainEntity.Customer, error)
	GetById(ctx context.Context, ID int64) (*customerDomainEntity.Customer, error)
	Delete(ctx context.Context, ID int64) error
	Update(ctx context.Context, ID int64, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error)
	Create(ctx context.Context, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error)
}
