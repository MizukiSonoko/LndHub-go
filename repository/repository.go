package repository

import e "github.com/MizukiSonoko/lnd-gateway/entity"

type UserRepo interface {
	Get(id string) e.User
	Update(e.User) error
	FindByPaymentHash(hash string) (e.User, error)
}

type userRepo struct {
}

func (repo *userRepo) Get(id string) e.User {
	return e.User{}
}

func (repo *userRepo) Update(e.User) error {
	return nil
}

func (repo *userRepo) FindByPaymentHash(hash string) (e.User, error) {
	return e.User{}, nil
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}
