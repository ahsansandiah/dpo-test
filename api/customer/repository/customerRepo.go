package customerRepository

import (
	"context"
	"database/sql"
	"time"

	customerDomainInterface "github.com/ahsansandiah/dpo-test/api/customer/domain"
	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
)

type Customer struct {
	DB  *sql.DB
	log log.Log
	cfg *config.Config
}

func NewCustomerRepository(mgr manager.Manager) customerDomainInterface.CustomerRepository {
	repo := new(Customer)
	repo.DB = mgr.GetDB()
	repo.log = mgr.GetLog()
	repo.cfg = mgr.GetConfig()

	return repo
}

func (r *Customer) GetAll(ctx context.Context, filter *customerDomainEntity.CustomerFilter) ([]customerDomainEntity.Customer, error) {
	query := "SELECT id, full_name, address, phone_number, email, is_active, created_at, updated_at FROM customers WHERE deleted_at IS NULL"
	var args []interface{}

	if filter.FullName != "" {
		query += " AND full_name = ?"
		args = append(args, filter.FullName)
	}

	if filter.PhoneNumber != "" {
		query += " AND phone_number = ?"
		args = append(args, filter.PhoneNumber)
	}

	if filter.Email != "" {
		query += " AND email = ?"
		args = append(args, filter.Email)
	}

	if filter.FullName == "" && filter.Email == "" && filter.PhoneNumber == "" {
		query = "SELECT id, full_name, address, phone_number, email, is_active, created_at, updated_at FROM customers WHERE deleted_at IS NULL"
	}

	query += " LIMIT ?"
	args = append(args, filter.LIMIT)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}
	defer rows.Close()

	var customers []customerDomainEntity.Customer
	for rows.Next() {
		var customer customerDomainEntity.Customer
		if err := rows.Scan(&customer.ID, &customer.FullName, &customer.Address, &customer.PhoneNumber, &customer.Email, &customer.IsActive, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			r.log.ErrorLog(ctx, err)
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return customers, nil
}

func (r *Customer) GetById(ctx context.Context, ID int64) (*customerDomainEntity.Customer, error) {
	customer := customerDomainEntity.Customer{}

	query := "SELECT id, full_name, address, phone_number, email, is_active, created_at, updated_at FROM customers WHERE id = ?"
	err := r.DB.QueryRow(query, ID).Scan(&customer.ID, &customer.FullName, &customer.Address, &customer.PhoneNumber, &customer.Email, &customer.IsActive, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &customer, nil
}

func (r *Customer) Delete(ctx context.Context, ID int64) error {
	stmt, err := r.DB.Prepare("UPDATE customers SET deleted_at = ? WHERE id = ?")
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

func (r *Customer) Update(ctx context.Context, ID int64, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error) {
	var customer customerDomainEntity.Customer
	stmt, err := r.DB.PrepareContext(ctx, "UPDATE customers SET full_name = ?, address = ?, phone_number = ?, email = ? WHERE id = ?")
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}
	defer stmt.Close()

	// Execute UPDATE statement
	_, err = stmt.ExecContext(ctx, request.FullName, request.Address, request.PhoneNumber, request.Email, ID)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	// Query the updated customer
	query := "SELECT id, full_name, address, phone_number, email, is_active, created_at, updated_at FROM customers WHERE id = ?"
	err = r.DB.QueryRowContext(ctx, query, ID).Scan(&customer.ID, &customer.FullName, &customer.Address, &customer.PhoneNumber, &customer.Email, &customer.IsActive, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &customer, nil
}

func (r *Customer) Create(ctx context.Context, request *customerDomainEntity.CustomerRequest) (*customerDomainEntity.Customer, error) {
	result, err := r.DB.ExecContext(ctx, "INSERT INTO customers (full_name, address, phone_number, email) VALUES (?, ?, ?, ?)", request.FullName, request.Address, request.PhoneNumber, request.Email)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	// Get ID of the inserted record
	customerID, err := result.LastInsertId()
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	var customer customerDomainEntity.Customer
	query := "SELECT id, full_name, address, phone_number, email, is_active, created_at, updated_at FROM customers WHERE id = ?"
	err = r.DB.QueryRowContext(ctx, query, customerID).Scan(&customer.ID, &customer.FullName, &customer.Address, &customer.PhoneNumber, &customer.Email, &customer.IsActive, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &customer, nil
}
