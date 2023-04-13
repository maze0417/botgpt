package mysql

import (
	"botgpt/internal/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

func GetMysqlDB() (*gorm.DB, error) {
	var errs error
	once.Do(func() {
		dsn := GetMysqlDsn()
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

		if err != nil {
			errs = err
			log.Errorf("Error connecting to database : error=%v", err)
		}

	})
	return db, errs
}

func GetMysqlDsn() string {
	c := config.GetConfig()
	DB_HOST := c.GetString("mysql.host")
	DB_PORT := c.GetString("mysql.port")
	DB_USERNAME := c.GetString("mysql.user")
	DB_PASSWORD := c.GetString("mysql.password")
	DB_NAME := c.GetString("mysql.database")

	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)
	return dsn
}

func CreateDbIfNotExist() {
	c := config.GetConfig()
	DB_HOST := c.GetString("mysql.host")
	DB_PORT := c.GetString("mysql.port")
	DB_USERNAME := c.GetString("mysql.user")
	DB_PASSWORD := c.GetString("mysql.password")
	DB_NAME := c.GetString("mysql.database")

	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/mysql?" + "parseTime=true&loc=Local"

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", DB_NAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatalf("Error connecting to database : error=%v", err)

	}
	tx := db.Exec(sql)
	if tx.Error != nil {
		log.Fatalln(tx.Error)
	}
}
