package repository

import (
	e "github.com/MizukiSonoko/LndHub-go/entity"
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func init() {
	d, err := leveldb.OpenFile("~/.lndhub-go/db", nil)
	if err != nil {
		panic(err)
	}
	db = d
}

type UserRepo interface {
	Get(id string) e.User
	Update(e.User) error
	FindByPaymentHash(hash string) (e.User, error)
}

type userRepo struct {
	db *leveldb.DB
}

func (repo *userRepo) Get(id string) e.User {
	panic("implement me")
}

func (repo *userRepo) GetBitcoinAddress(id string) (string, error) {
	address, err := repo.db.Get([]byte("bitcoin_address_for_"+id), nil)
	if err != nil {
		return "", err
	}
	return string(address), nil
}

func (repo *userRepo) SetBitcoinAddress(id, address string) error {
	return db.Put([]byte("bitcoin_address_for_"+id), []byte("address"), nil)
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
