package database

import (
	"fmt"

	"github.com/mesh-dell/todo-list-API/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDB struct {
	DbClient *gorm.DB
}

func NewGormDb(config config.Config) *GormDB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPassword,
		config.DbAddr,
		config.DbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}
	return &GormDB{
		DbClient: db,
	}
}
