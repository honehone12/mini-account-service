package db

import (
	"errors"
	"log"
	"mini-account-service/db/controller"
	"mini-account-service/db/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Orm struct {
	db *gorm.DB
}

func (orm Orm) User() controller.UserController {
	return controller.NewUserController(orm.db)
}

func NewOrm() (Orm, error) {
	orm := Orm{}
	var err error

	dsn := os.Getenv("POSTGRES_DSN")
	if len(dsn) == 0 {
		return orm, errors.New("env param POSTGRES_DSN is empty")
	}

	orm.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return orm, err
	}

	if err = orm.migrate(); err != nil {
		return orm, err
	}

	log.Println("new database connection is done")
	return orm, nil
}

func (conn Orm) migrate() error {
	if err := conn.db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}
