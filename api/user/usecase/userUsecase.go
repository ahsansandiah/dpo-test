package userUsecase

import (
	"context"
	"errors"

	userDomainInterface "github.com/ahsansandiah/dpo-test/api/user/domain"
	userDomainEntity "github.com/ahsansandiah/dpo-test/api/user/domain/entity"
	userRepository "github.com/ahsansandiah/dpo-test/api/user/repository"
	jwtAuth "github.com/ahsansandiah/dpo-test/packages/auth/jwt"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	log  log.Log
	cfg  *config.Config
	repo userDomainInterface.UserRepository
	jwt  jwtAuth.Jwt
}

func NewUserUsecase(mgr manager.Manager) userDomainInterface.UserUsecase {
	usecase := new(UserUsecase)
	usecase.log = mgr.GetLog()
	usecase.cfg = mgr.GetConfig()
	usecase.repo = userRepository.NewUserRepository(mgr)
	usecase.jwt = mgr.GetJwt()

	return usecase
}

func (u *UserUsecase) GetUserLogin(ctx context.Context, token string) (*userDomainEntity.UserResponse, error) {
	jwtData, err := u.jwt.ExtractJwtToken(token)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Failed get user detail")
		return nil, errMsg
	}

	customer, err := u.repo.GetById(ctx, jwtData.UserID)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error fetching customer details")
		return nil, errMsg
	}

	result := &userDomainEntity.UserResponse{
		ID:        customer.ID,
		Username:  customer.Username,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	return result, nil
}

func (u *UserUsecase) Create(ctx context.Context, request *userDomainEntity.UserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error inserting user")
		return errMsg
	}

	request.PasswordHash = hashedPassword
	err = u.repo.Create(ctx, request)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Error inserting user")
		return errMsg
	}

	return nil
}

func (u *UserUsecase) Login(ctx context.Context, request *userDomainEntity.LoginRequest) (*userDomainEntity.LoginResponse, error) {
	// get user by username
	user, err := u.repo.GetByUsername(ctx, request.Username)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Username does not match")
		return nil, errMsg
	}

	// Example of verifying a hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Password does not match")
		return nil, errMsg
	}

	dataJwt := &jwtAuth.JwtData{
		UserID:    int64(user.ID),
		Reference: user.Username,
	}

	// generate token
	accessToken, expiredTime, err := u.jwt.GenerateToken(dataJwt)
	if err != nil {
		u.log.ErrorLog(ctx, err)
		errMsg := errors.New("Failed to login")
		return nil, errMsg
	}

	result := &userDomainEntity.LoginResponse{
		Token:      accessToken,
		ExpiryTime: expiredTime,
	}

	return result, nil
}

func (u *UserUsecase) Logout(ctx context.Context) error {

	return nil
}
