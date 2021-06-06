package repository

import (
	"github.com/jinzhu/gorm"
	"user/domain/model"
)

type IUserRepository interface {
	InitTable() error
	FindUserByName(string) (*model.User, error)
	FindUserByID(int64) (*model.User, error)
	CreateUser(user *model.User) (int64, error)
	DeleteUserByID(int64) error
	UpdateUser(user *model.User) error
	FindAll() ([]model.User, error)
}

type UserRepository struct {
	mysqlDB *gorm.DB
}

func (u *UserRepository) InitTable() error {
	return u.mysqlDB.CreateTable(&model.User{}).Error
}

func (u *UserRepository) FindUserByName(username string) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqlDB.Where("user_name = ?", username).Find(user).Error
}

func (u *UserRepository) FindUserByID(id int64) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqlDB.Where("id = ?", id).Find(user).Error

}

func (u *UserRepository) CreateUser(user *model.User) (int64, error) {
	return user.ID, u.mysqlDB.Create(user).Error

}

func (u *UserRepository) DeleteUserByID(userID int64) error {
	return u.mysqlDB.Where("id = ?", userID).Delete(&model.User{}).Error
}

func (u *UserRepository) UpdateUser(user *model.User) error {
	return u.mysqlDB.Model(user).Update(&user).Error
}

func (u *UserRepository) FindAll() (userAll []model.User, err error) {
	return userAll, u.mysqlDB.Find(&userAll).Error

}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{mysqlDB: db}

}
