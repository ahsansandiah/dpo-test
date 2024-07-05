package userHandler

import (
	"encoding/json"
	"net/http"

	userDomainInterface "github.com/ahsansandiah/dpo-test/api/user/domain"
	userDomainEntity "github.com/ahsansandiah/dpo-test/api/user/domain/entity"
	userUsecase "github.com/ahsansandiah/dpo-test/api/user/usecase"
	jwtAuth "github.com/ahsansandiah/dpo-test/packages/auth/jwt"
	middlewareAuth "github.com/ahsansandiah/dpo-test/packages/auth/middleware"
	res "github.com/ahsansandiah/dpo-test/packages/json"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
)

type User struct {
	log        log.Log
	Json       res.Json
	Usecase    userDomainInterface.UserUsecase
	jwt        jwtAuth.Jwt
	Middleware middlewareAuth.Middleware
}

func NewUserHandler(mgr manager.Manager) userDomainInterface.UserHandler {
	handler := new(User)
	handler.Usecase = userUsecase.NewUserUsecase(mgr)
	handler.Json = mgr.GetJson()
	handler.Middleware = mgr.GetMiddleware()

	return handler
}

func (h *User) Detail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := h.Middleware.GetTokenInHeader(r)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		customer, err := h.Usecase.GetUserLogin(ctx, token)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success get data", customer)
	})
}

func (h *User) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req *userDomainEntity.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := req.Validate(); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err := h.Usecase.Create(ctx, req)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success created", nil)
	})
}

func (h *User) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req *userDomainEntity.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := req.LoginValidate(); err != nil {
			h.Json.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		result, err := h.Usecase.Login(ctx, req)
		if err != nil {
			h.Json.ErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		h.Json.SuccessResponse(w, r, http.StatusCreated, "Success created", result)
	})
}
