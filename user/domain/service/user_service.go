package service

import (
	"golang.org/x/crypto/bcrypt"
	"user/domain/model"
	"user/domain/repository"
)

type IUserDataService interface {
	AddUser(user *model.User) (int64, error)
	DeleteUser(int64) error
	UpdateUser(user *model.User, isChangePwd bool) error
	FindUserByName(string) (*model.User, error)
	CheckPwd(userName string, pwd string) (isOk bool, err error)
}

type UserDataService struct {
	UserRepository repository.IUserRepository
}

// 加密
func GeneratePassword(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

// 解密
func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}

func (u UserDataService) AddUser(user *model.User) (int64, error) {
	pwdByte, err := GeneratePassword(user.Pwd)
	if err != nil {
		return user.ID, err
	}
	user.Pwd = string(pwdByte)
	return u.UserRepository.CreateUser(user)
}

func (u UserDataService) DeleteUser(userID int64) error {
	return u.UserRepository.DeleteUserByID(userID)

}

func (u UserDataService) UpdateUser(user *model.User, isChangePwd bool) error {
	if isChangePwd {
		pwdByte, err := GeneratePassword(user.Pwd)
		if err != nil {
			return err
		}
		user.Pwd = string(pwdByte)
	}
	return u.UserRepository.UpdateUser(user)
}

func (u UserDataService) FindUserByName(userName string) (*model.User, error) {
	return u.UserRepository.FindUserByName(userName)
}

func (u UserDataService) CheckPwd(userName string, pwd string) (isOk bool, err error) {
	user, err := u.UserRepository.FindUserByName(userName)
	if err != nil {
		return false, err
	}
	return ValidatePassword(pwd, user.Pwd)
}

func NewUserDataService(userRepository repository.IUserRepository) IUserDataService {
	return &UserDataService{UserRepository: userRepository}
}
