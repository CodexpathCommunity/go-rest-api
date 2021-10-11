package user

import (
	"context"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Repository encapsulates the logic to access users from the data source.
type UsersRepository interface {
	GetUser(ctx context.Context, id string) (entity.Users, error)
	CreateUser(ctx context.Context, album entity.Users) error
	UpdateUser(ctx context.Context, album entity.Users) error
	DeleteUser(ctx context.Context, id string) error
}

type usersRepository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewUsersRepository(db *dbcontext.DB, logger log.Logger) UsersRepository {
	return usersRepository{db, logger}
}

func (r usersRepository) GetUser(ctx context.Context, ID string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Model(ID, &user)
	return user, err
}

func (r usersRepository) CreateUser(ctx context.Context, user entity.Users) error {
	return r.db.With(ctx).Model(&user).Insert()
}

func (r usersRepository) UpdateUser(ctx context.Context, user entity.Users) error {
	return r.db.With(ctx).Model(&user).Update()
}

func (r usersRepository) DeleteUser(ctx context.Context, emailAddress string) error {
	user, err := r.GetUser(ctx, emailAddress)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&user).Delete()
}