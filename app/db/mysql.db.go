package db

import (
	"fmt"

	"de.server/app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToMySQLDb(configs model.Conf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", configs.User, configs.Password, configs.Host, configs.Port, configs.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
