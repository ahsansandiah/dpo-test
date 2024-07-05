package userDomainInterface

import (
	"context"
	"net/http"

	userDomainEntity "github.com/ahsansandiah/dpo-test/api/user/domain/entity"
)

type UserHandler interface {
	Detail() http.Handler
	Create() http.Handler
	Login() http.Handler
}

type UserUsecase interface {
	GetUserLogin(ctx context.Context, token string) (*userDomainEntity.UserResponse, error)
	Create(ctx context.Context, request *userDomainEntity.UserRequest) error
	Login(ctx context.Context, request *userDomainEntity.LoginRequest) (*userDomainEntity.LoginResponse, error)
}

type UserRepository interface {
	GetById(ctx context.Context, ID int64) (*userDomainEntity.User, error)
	Create(ctx context.Context, request *userDomainEntity.UserRequest) error
	GetByUsername(ctx context.Context, username string) (*userDomainEntity.User, error)
}
