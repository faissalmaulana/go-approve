package db

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/faissalmaulana/go-approve/internal/model"
	_ "github.com/joho/godotenv/autoload"
)

var (
	dbname   = os.Getenv("DB_DATABASE_NAME")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	db       *gorm.DB
)

func New(log *zap.Logger) *gorm.DB {
	// reuse connection
	if db != nil {
		return db
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error", zap.Error(err))
	}

	db.AutoMigrate(&model.User{})

	log.Info("DB connected", zap.String("addr", fmt.Sprintf("%s:%s", host, port)))

	return db
}
