package internal

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type DatabaseService struct {
	Config DatabaseConfig
	Db     *gorm.DB
}

func NewDatabaseService(config *DatabaseConfig) *DatabaseService {
	return &DatabaseService{
		Db: connect(config),
	}
}

func connect(config *DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open("host="+config.Host+" user="+config.Username+" password="+config.Password+" port="+config.Port), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err.Error())
	}
	return db
}
