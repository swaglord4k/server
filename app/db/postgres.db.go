package db

import (
	"fmt"

	"de.server/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPostgresDb(configs model.Conf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require", configs.Host, configs.User, configs.Password, configs.Dbname, configs.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
