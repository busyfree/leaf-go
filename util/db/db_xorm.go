package db

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"

	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/log"
)

var dbXORMs = make(map[string]*xorm.Engine, 4)

func GetXORM(ctx context.Context, name string) *xorm.Engine {
	lock.RLock()
	db := dbXORMs[name]
	lock.RUnlock()

	if db != nil {
		return db
	}

	dsn := conf.GetString("DB_" + name + "_DSN")
	var logger = log.Get(ctx)
	// logger.Info("dsn:", dsn)
	sqldb, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Get(ctx).Panic(err)
	}
	maxCon := conf.GetInt("DB_" + name + "_MAX_CONN")
	maxIdle := conf.GetInt("DB_" + name + "_MAX_IDLE")
	maxLife := conf.GetInt("DB_" + name + "_MAX_LIFE")

	if maxCon <= 0 {
		maxCon = 10
	}
	if maxIdle <= 0 {
		maxIdle = 10
	}
	if maxLife <= 0 {
		maxLife = 60
	}
	sqldb.SetMaxOpenConns(maxCon)
	sqldb.SetMaxIdleConns(maxIdle)
	sqldb.SetConnMaxLifetime(time.Duration(maxLife) * time.Second)
	// logger.Infof("sqldb:%v", sqldb)
	errping := sqldb.Ping()
	logger.Infof("dbPingErr:%v", errping)
	lock.Lock()
	dbXORMs[name] = sqldb
	lock.Unlock()
	return sqldb
}

// ResetXORM 关闭所有 DB 连接
// 新调用 GetXORM 方法时会使用最新 DB 配置创建连接
func ResetXORM() {
	for k, db := range dbXORMs {
		db.Close()
		delete(dbXORMs, k)
	}
}
