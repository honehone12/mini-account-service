package controller

import (
	"errors"
	"mini-account-service/db/models"
	"mini-account-service/db/pwhasher"

	"gorm.io/gorm"
)

var ErrorInvalidPassword = errors.New("invalid password")

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return UserController{db}
}

func (c UserController) Create(email string, password string) error {
	user, err := models.NewUser(email, password)
	if err != nil {
		return err
	}

	result := c.db.Create(user)
	return result.Error
}

func (c UserController) ReadByEmail(email string) (*models.User, error) {
	user := &models.User{}
	result := c.db.Select(
		"ID", "CreatedAt", "UpdatedAt", "Uuid", "Email",
	).Where("email = ?", email).Take(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (c UserController) ReadByUuid(uuid string) (*models.User, error) {
	user := &models.User{}
	result := c.db.Select(
		"ID", "CreatedAt", "UpdatedAt", "Uuid", "Email",
	).Where("uuid = ?", uuid).Take(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (c UserController) VerifyPassword(email string, password string) (string, error) {
	user := &models.User{}
	result := c.db.Select("PasswordHash", "Salt", "Uuid").Where("email = ?", email).Take(user)
	if result.Error != nil {
		return "", result.Error
	}

	hasher := pwhasher.NewPasswordHasher(password)
	ok, err := hasher.Verify(user.PasswordHash, user.Salt)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrorInvalidPassword
	}

	return user.Uuid, nil
}
