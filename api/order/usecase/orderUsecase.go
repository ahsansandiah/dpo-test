package orderUsecase

import (
	"context"
	"errors"

	orderDomainInterface "github.com/ahsansandiah/dpo-test/api/order/domain"
	orderDomainEntity "github.com/ahsansandiah/dpo-test/api/order/domain/entity"
	orderRepository "github.com/ahsansandiah/dpo-test/api/order/repository"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
)

type OrderUsecase struct {
	log  log.Log
	cfg  *config.Config
	repo orderDomainInterface.OrderRepository
}

func NewOrderUsecase(mgr manager.Manager) orderDomainInterface.OrderUsecase {
	usecase := new(OrderUsecase)
	usecase.log = mgr.GetLog()
	usecase.cfg = mgr.GetConfig()
	usecase.repo = orderRepository.NewOrderRepository(mgr)

	return usecase
}

func (u *OrderUsecase) GetAll(ctx context.Context, filter *orderDomainEntity.OrderFilter) ([]orderDomainEntity.OrderResponse, error) {
	orders, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		return nil, err
	}

	return orders, nil
}

func (u *OrderUsecase) Delete(ctx context.Context, ID int64) error {
	order, err := u.repo.GetById(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching order details")
		return errMsg
	}

	err = u.repo.Delete(ctx, order.ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error deleting order")
		return errMsg
	}

	return nil
}

func (u *OrderUsecase) GetByID(ctx context.Context, ID int64) (*orderDomainEntity.OrderResponse, error) {
	// get order
	order, err := u.repo.GetById(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching order")
		return nil, errMsg
	}

	customer, err := u.repo.GetCustomer(ctx, int64(order.CustomerID))
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching order")
		return nil, errMsg
	}

	// get order items
	orderItems, err := u.repo.GetOrderItems(ctx, order.ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching order items")
		return nil, errMsg
	}

	result := orderDomainEntity.OrderResponse{
		ID:          order.ID,
		OrderDate:   order.OrderDate,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Customer:    customer,
		Items:       orderItems,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	return &result, nil
}

func (u *OrderUsecase) Update(ctx context.Context, ID int64, request *orderDomainEntity.OrderUpdateRequest) (*orderDomainEntity.OrderResponse, error) {
	order, err := u.repo.GetById(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching order details")
		return nil, errMsg
	}

	if request.TotalAmount <= 0 {
		request.TotalAmount = order.TotalAmount
	}

	if request.OrderDate.IsZero() {
		request.OrderDate = order.OrderDate
	}

	_, err = u.repo.Update(ctx, ID, request)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error update order")
		return nil, errMsg
	}

	result, err := u.GetByID(ctx, order.ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching order details")
		return nil, errMsg
	}

	return result, nil
}

func (u *OrderUsecase) Create(ctx context.Context, request *orderDomainEntity.OrderRequest) (*orderDomainEntity.OrderRequest, error) {
	// check customer
	if !u.ValidateCustomer(ctx, request.CustomerID) {
		errMsg := errors.New("Error customer not found")
		return nil, errMsg
	}

	// create order with order items
	err := u.repo.Create(ctx, request)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error inserting order")
		return nil, errMsg
	}

	return request, nil
}

func (u *OrderUsecase) ValidateCustomer(ctx context.Context, customerID int64) bool {
	_, err := u.repo.GetCustomer(ctx, customerID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		return false
	}

	return true
}
