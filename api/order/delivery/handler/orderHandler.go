package orderHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	orderDomainInterface "github.com/ahsansandiah/dpo-test/api/order/domain"
	orderDomainEntity "github.com/ahsansandiah/dpo-test/api/order/domain/entity"
	orderUsecase "github.com/ahsansandiah/dpo-test/api/order/usecase"
	res "github.com/ahsansandiah/dpo-test/packages/json"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

type Order struct {
	log     log.Log
	Json    res.Json
	Usecase orderDomainInterface.OrderUsecase
}

func NewOrderHandler(mgr manager.Manager) orderDomainInterface.OrderHandler {
	handler := new(Order)
	handler.Usecase = orderUsecase.NewOrderUsecase(mgr)
	handler.Json = mgr.GetJson()

	return handler
}

func (h *Order) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		queryParams := r.URL.Query()

		limitStr := queryParams.Get("limit")
		limit := 10
		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err == nil {
				limit = l
			}
		}

		filter := &orderDomainEntity.OrderFilter{
			OrderDate:  queryParams.Get("order_date"),
			CustomerID: queryParams.Get("customer_id"),
			Status:     queryParams.Get("status"),
			LIMIT:      limit,
		}

		result, err := h.Usecase.GetAll(ctx, filter)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success get data", result)
	})
}

func (h *Order) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIDStr := mux.Vars(r)["id"]
		orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		err = h.Usecase.Delete(ctx, orderID)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, fmt.Sprintf("Order with ID %d deleted successfully", orderID), nil)
	})
}

func (h *Order) GetByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIDStr := mux.Vars(r)["id"]
		orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		order, err := h.Usecase.GetByID(ctx, orderID)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success get data", order)
	})
}

func (h *Order) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIDStr := mux.Vars(r)["id"]
		orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		var req *orderDomainEntity.OrderUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		order, err := h.Usecase.Update(ctx, orderID, req)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success updated", order)
	})
}

func (h *Order) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req *orderDomainEntity.OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := req.Validate(); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		order, err := h.Usecase.Create(ctx, req)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success created", order)
	})
}
