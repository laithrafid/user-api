package services

import (
	"github.com/laithrafid/user-api/src/domain/users"
	"github.com/laithrafid/utils-go/crypto_utils"
	"github.com/laithrafid/utils-go/date_utils"
	"github.com/laithrafid/utils-go/errors_utils"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, errors_utils.RestErr)
	CreateUser(users.User) (*users.User, errors_utils.RestErr)
	UpdateUser(bool, users.User) (*users.User, errors_utils.RestErr)
	DeleteUser(int64) errors_utils.RestErr
	SearchUser(string) (users.Users, errors_utils.RestErr)
	LoginUser(users.LoginRequest) (*users.User, errors_utils.RestErr)
}

func (s *usersService) GetUser(userId int64) (*users.User, errors_utils.RestErr) {
	dao := &users.User{Id: userId}
	if err := dao.Get(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, errors_utils.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, errors_utils.RestErr) {
	current := &users.User{Id: user.Id}
	if err := current.Get(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}

		if user.LastName != "" {
			current.LastName = user.LastName
		}

		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId int64) errors_utils.RestErr {
	dao := &users.User{Id: userId}
	return dao.Delete()
}

func (s *usersService) SearchUser(status string) (users.Users, errors_utils.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, errors_utils.RestErr) {
	dao := &users.User{
		Email:    request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
