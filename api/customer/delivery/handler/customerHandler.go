package customerHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	customerDomainInterface "github.com/ahsansandiah/dpo-test/api/customer/domain"
	customerDomainEntity "github.com/ahsansandiah/dpo-test/api/customer/domain/entity"
	customerUsecase "github.com/ahsansandiah/dpo-test/api/customer/usecase"
	res "github.com/ahsansandiah/dpo-test/packages/json"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

type Customer struct {
	log     log.Log
	Json    res.Json
	Usecase customerDomainInterface.CustomerUsecase
}

func NewCustomerHandler(mgr manager.Manager) customerDomainInterface.CustomerHandler {
	handler := new(Customer)
	handler.Usecase = customerUsecase.NewCustomerUsecase(mgr)
	handler.Json = mgr.GetJson()

	return handler
}

func (h *Customer) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		queryParams := r.URL.Query()
		limit := queryParams.Get("limit")
		if limit == "" {
			limit = "10"
		}

		isActiveParams := queryParams.Get("is_active")
		isActive := false
		if isActiveParams == "true" {
			isActive = true
		}

		filter := &customerDomainEntity.CustomerFilter{
			FullName:    queryParams.Get("full_name"),
			Email:       queryParams.Get("email"),
			PhoneNumber: queryParams.Get("phone_number"),
			IsActive:    isActive,
			LIMIT:       limit,
		}

		result, err := h.Usecase.GetAll(ctx, filter)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success get data", result)
	})
}

func (h *Customer) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerIDStr := mux.Vars(r)["id"]
		customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		err = h.Usecase.Delete(ctx, customerID)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, fmt.Sprintf("Customer with ID %d deleted successfully", customerID), nil)
	})
}

func (h *Customer) GetByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerIDStr := mux.Vars(r)["id"]
		customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		customer, err := h.Usecase.GetByID(ctx, customerID)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success get data", customer)
	})
}

func (h *Customer) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerIDStr := mux.Vars(r)["id"]
		customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		var req *customerDomainEntity.CustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		customer, err := h.Usecase.Update(ctx, customerID, req)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success updated", customer)
	})
}

func (h *Customer) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req *customerDomainEntity.CustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := req.Validate(); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		customer, err := h.Usecase.Create(ctx, req)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success created", customer)
	})
}
