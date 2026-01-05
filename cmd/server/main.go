package main

import (
	"database/sql"
	"log"

	"minigo/internal/infrastructure/config"
	"minigo/internal/infrastructure/id"
	"minigo/internal/infrastructure/logging"
	httpx "minigo/internal/interfaces/http"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// initConfig 初始化配置和日志
func initConfig() {
	config.Init()
	logging.Init(config.GetLogLevel())
}

// connectDB 连接数据库并返回bun.DB实例
func connectDB() (*bun.DB, error) {
	dsn := viper.GetString("DB_DSN")
	if dsn == "" {
		dsn = config.GetDBDsn()
	}

	hsqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(hsqldb, pgdialect.New())
	if config.IsDevEnv() {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, nil
}

func main() {
	initConfig()
	id.Init()

	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	r := httpx.BuildRouter(db)

	port := viper.GetString("PORT")
	if port == "" {
		port = config.GetPort()
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
