package mysql

import (
	"fmt"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	rwDB *gorm.DB
	roDB *gorm.DB
)

const (
	DBMaxOpenConn = 10
	DBMaxIdleConn = 5
	DBMaxLifeTime = time.Second * 30
)

func Connect() (*gorm.DB, *gorm.DB) {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("mysql.NewConfig err: %w", err))
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             300 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsnRW := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true", config.MySQLUser, config.MySQLPass, config.MySQLHostRW, config.MySQLDBName)
	rwDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                     dsnRW,
		DontSupportRenameColumn: false,
		DontSupportRenameIndex:  false,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger,
		PrepareStmt:                              true,
	})
	if err != nil {
		panic(fmt.Errorf("mysql rw: %s", err))
	}

	rwSQL, err := rwDB.DB()
	if err != nil {
		panic(fmt.Errorf("rwDB.DB() err: %w", err))
	}
	rwSQL.SetMaxOpenConns(DBMaxOpenConn)
	rwSQL.SetMaxIdleConns(DBMaxIdleConn)
	rwSQL.SetConnMaxLifetime(DBMaxLifeTime)
	err = rwSQL.Ping()
	if err != nil {
		panic(fmt.Errorf("mysql rw: %s", err))
	}
	err = rwDB.AutoMigrate(&domain.Setting{}, &domain.TVOrder{})
	if err != nil {
		panic(fmt.Errorf("AutoMigrate failed err: %w", err))
	}

	dsnRO := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true", config.MySQLUser, config.MySQLPass, config.MySQLHostRO, config.MySQLDBName)
	roDB, err = gorm.Open(mysql.Open(dsnRO), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger,
		PrepareStmt:                              true,
	})
	if err != nil {
		panic(fmt.Errorf("mysql ro: %s", err))
	}

	roSQL, err := roDB.DB()
	if err != nil {
		panic(fmt.Errorf("roDB.DB() err: %w", err))
	}
	roSQL.SetMaxOpenConns(DBMaxOpenConn)
	roSQL.SetMaxIdleConns(DBMaxIdleConn)
	roSQL.SetConnMaxLifetime(DBMaxLifeTime)
	err = roSQL.Ping()
	if err != nil {
		panic(fmt.Errorf("mysql ro: %s", err))
	}

	return rwDB, roDB
}

func CloseConn() {
	rwSQL, _ := rwDB.DB()
	rwSQL.Close()

	roSQL, _ := roDB.DB()
	roSQL.Close()
}
