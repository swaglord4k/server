package db

import (
	"fmt"

	m "de.server/app/model"
	"gorm.io/gorm"
)

func ConnectToDB(confType m.ConfigType) (*gorm.DB, error) {
	var configs m.Conf
	configs.GetConf(confType)
	switch confType {
	case m.MysqlConfig:
		return ConnectToMySQLDb(configs)
	case m.PostgresConfig:
		return ConnectToPostgresDb(configs)
	}
	return nil, fmt.Errorf("config type %s not supported", confType)
}
