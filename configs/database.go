package configs

import (
	"MyGram/app"
	"MyGram/models"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func ConfigDatabase() {
	var logLevel gormlogger.LogLevel
	switch viper.GetString("LOG_DB_LEVEL") {
	case "SILENT":
		logLevel = gormlogger.Silent
	case "ERROR":
		logLevel = gormlogger.Error
	case "WARN":
		logLevel = gormlogger.Warn

	case "INFO":
		logLevel = gormlogger.Info
	}
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Duration(viper.GetInt64("LOG_DB_SLOWQUERY_THRESHOLD")) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logLevel,                                                                       // Log level
			IgnoreRecordNotFoundError: true,                                                                           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,
			// Disable color
		},
	)
	//&parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&collation=utf8mb4_unicode_ci&loc=%s",
		viper.GetString("DB_USERNAME"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
		url.QueryEscape(viper.GetString("DB_TIMEZONE")))
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true, Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	app.Db = db
	app.Db.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})
}
