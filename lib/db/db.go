package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	utility "github.com/kyken/fiber-sandbox/lib/utils"
)

// データベース構造体
type Database struct {
	*sqlx.DB
}

// データベース接続設定
func NewDatabase() (*Database, error) {
	// 環境変数から設定を読み込み
	dbHost := utility.GetEnv("DB_HOST", "127.0.0.1")
	dbPort := utility.GetEnv("DB_PORT", "3306")
	dbUser := utility.GetEnv("DB_USER", "appuser")
	dbPass := utility.GetEnv("DB_PASSWORD", "apppassword")
	dbName := utility.GetEnv("DB_NAME", "appdb")

	// DSN（Data Source Name）の構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// データベース接続
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 接続設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{db}, nil
}
