package dao

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"yuguosheng/int/mychatops/middleware"
)

const (
	DB_NAME = "/db/data.db"
)

var Db *gorm.DB

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Context struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name"`
	Command     string    `gorm:"column:command"`
	Context_msg string    `gorm:"column:context_msg"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type ContextCommand struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Text      string    `gorm:"column:text"`
	Command   string    `gorm:"column:command"`
	Auth      string    `gorm:"column:auth"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func LoadDatabase() {
	middleware.MyLogger.Info("初始化Db.....")
	WorkPath, _ := os.Getwd()
	// 打开数据库

	db, err := gorm.Open(sqlite.Open(WorkPath+DB_NAME), &gorm.Config{})
	if err != nil {
		middleware.MyLogger.Fatal("初始化DB打开数据库失败", zap.Any("Error", err))
	}

	// 创建表
	err = db.AutoMigrate(&User{})
	if err != nil {
		middleware.MyLogger.Fatal("初始化DB: 创建表失败", zap.Any("error", err))
	}

	err = db.AutoMigrate(&Context{})
	if err != nil {
		middleware.MyLogger.Fatal("初始化DB: 创建表失败", zap.Any("error", err))
	}

	err = db.AutoMigrate(&ContextCommand{})
	if err != nil {
		middleware.MyLogger.Fatal("初始化DB: 创建表失败", zap.Any("error", err))
	}

	Db = db
}
