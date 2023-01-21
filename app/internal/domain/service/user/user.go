package user

import (
	"avito-tech/app/adapters/db/postgresql"
	"avito-tech/app/internal/domain/entity"
	"context"
)

type Storage interface {
	Create(ctx context.Context) (int, error)
	FindOne(ctx context.Context, id int) (*entity.User, error)
	ChargeBalance(ctx context.Context, user *entity.User) error
	RemoveBalance(ctx context.Context, user *entity.User) error
	TransferMoney(ctx context.Context, data *entity.TransferMoney) error
}

type UserService struct {
	storage *postgresql.UserStorage
}

func NewUserService(storage *postgresql.UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (u UserService) CreateUser(ctx context.Context) (int, error) {
	return u.storage.Create(ctx)
}

func (u UserService) FindUserById(ctx context.Context, id int) (*entity.User, error) {
	return u.storage.FindOne(ctx, id)
}

func (u UserService) ChargeBalance(ctx context.Context, userChargeMoneyDTO UserChargeMoneyDTO) error {
	user := &entity.User{
		Id:      userChargeMoneyDTO.Id,
		Balance: userChargeMoneyDTO.Amount,
	}

	return u.storage.ChargeBalance(ctx, user)
}

func (u UserService) RemoveBalance(ctx context.Context, userRemoveMoneyDTO UserRemoveMoneyDTO) error {
	user := &entity.User{
		Id:      userRemoveMoneyDTO.Id,
		Balance: userRemoveMoneyDTO.Amount,
	}

	return u.storage.RemoveBalance(ctx, user)
}

func (u UserService) TransferMoney(ctx context.Context, data TransferMoneyDTO) error {
	dataEntity := &entity.TransferMoney{
		From:   data.From,
		To:     data.To,
		Amount: data.Amount,
	}

	return u.storage.TransferMoney(ctx, dataEntity)
}
