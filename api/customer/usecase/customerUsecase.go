package customerUsecase

import (
	"context"
	"errors"

	customerDomainInterface "github.com/ahsansandiah/dpo-test/api/customer/domain"
	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
	customerRepository "github.com/ahsansandiah/dpo-test/api/customer/repository"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
)

type CustomerUsecase struct {
	log  log.Log
	cfg  *config.Config
	repo customerDomainInterface.CustomerRepository
}

func NewCustomerUsecase(mgr manager.Manager) customerDomainInterface.CustomerUsecase {
	usecase := new(CustomerUsecase)
	usecase.log = mgr.GetLog()
	usecase.cfg = mgr.GetConfig()
	usecase.repo = customerRepository.NewCustomerRepository(mgr)

	return usecase
}

func (u *CustomerUsecase) GetAll(ctx context.Context, filter *customerDomainEntity.CustomerFilter) ([]customerDomainEntity.Customer, error) {
	customers, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		return nil, err
	}

	return customers, nil
}

func (u *CustomerUsecase) Delete(ctx context.Context, ID int64) error {
	_, err := u.repo.GetById(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching customer details")
		return errMsg
	}

	err = u.repo.Delete(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error deleting customer")
		return errMsg
	}

	return nil
}

func (u *CustomerUsecase) GetByID(ctx context.Context, ID int64) (*customerDomainEntity.Customer, error) {
	customer, err := u.repo.GetById(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching customer details")
		return nil, errMsg
	}

	return customer, nil
}

func (u *CustomerUsecase) Update(ctx context.Context, ID int64, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error) {
	customer, err := u.repo.GetById(ctx, ID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching cake details")
		return nil, errMsg
	}

	if request.FullName == "" {
		request.FullName = customer.FullName
	}

	if request.Address == "" {
		request.Address = customer.Address
	}

	if request.PhoneNumber == "" {
		request.PhoneNumber = customer.PhoneNumber
	}

	if request.Email == "" {
		request.Email = customer.Email
	}

	result, err := u.repo.Update(ctx, ID, request)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error deleting customer")
		return nil, errMsg
	}

	return result, nil
}

func (u *CustomerUsecase) Create(ctx context.Context, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error) {
	custmer, err := u.repo.Create(ctx, request)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error inserting customer")
		return nil, errMsg
	}

	return custmer, nil
}
