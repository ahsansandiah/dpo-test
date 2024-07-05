package userRepository

import (
	"context"
	"database/sql"

	userDomainInterface "github.com/ahsansandiah/dpo-test/api/user/domain"
	userDomainEntity "github.com/ahsansandiah/dpo-test/api/user/domain/entity"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/manager"
)

type User struct {
	DB  *sql.DB
	log log.Log
	cfg *config.Config
}

func NewUserRepository(mgr manager.Manager) userDomainInterface.UserRepository {
	repo := new(User)
	repo.DB = mgr.GetDB()
	repo.log = mgr.GetLog()
	repo.cfg = mgr.GetConfig()

	return repo
}

func (r *User) GetById(ctx context.Context, ID int64) (*userDomainEntity.User, error) {
	user := userDomainEntity.User{}

	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?"
	err := r.DB.QueryRow(query, ID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &user, nil
}

func (r *User) Create(ctx context.Context, request *userDomainEntity.UserRequest) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)", request.Username, request.PasswordHash, request.Email)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return err
	}

	return nil
}

func (r *User) GetByUsername(ctx context.Context, username string) (*userDomainEntity.User, error) {
	user := userDomainEntity.User{}

	query := "SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = ?"
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.log.ErrorLog(ctx, err)
		return nil, err
	}

	return &user, nil
}
