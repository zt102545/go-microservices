package com

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-microservices/dao/mysql/model"
	redisModel "go-microservices/dao/redis/model"
	"log"
	"time"
)

type ComConfig struct {
	DB    *sql.DB
	Conn  sqlx.SqlConn
	Rc    *redis.Client
	Cache *collection.Cache

	// mysql model
	model.UserModel
	model.GlobalVariablesModel
	// redis model
	redisModel.LockRedis
}

type ComFunc struct {
	ComConfig
}

type Config struct {
	ComConfig
	ComFuncInterface
}

func InitCom(db *sql.DB, postgres map[string]*sql.DB, rc *redis.Client) *Config {

	connMap := map[string]sqlx.SqlConn{
		"mysql": sqlx.NewSqlConnFromDB(db),
	}

	for k, v := range postgres {
		connMap[k] = sqlx.NewSqlConnFromDB(v)
	}

	cache, err := collection.NewCache(5*time.Minute, collection.WithName("any"))
	if err != nil {
		log.Fatal(err)
	}
	com := &ComConfig{
		Conn:  connMap["mysql"],
		DB:    db,
		Rc:    rc,
		Cache: cache,

		// mysql model
		UserModel:            model.NewUserModel(connMap["mysql"]),
		GlobalVariablesModel: model.NewGlobalVariablesModel(connMap["mysql"]),
		//postgres model

		//redis model
		LockRedis: redisModel.NewLockRedis(rc),
	}
	return &Config{
		ComConfig: *com,
		ComFuncInterface: &ComFunc{
			ComConfig: *com,
		},
	}
}

// 公共方法接口
type ComFuncInterface interface {
	// CheckUser 用户校验
	CheckUser(user *model.User) bool
	GetGlobalVariable(ctx context.Context, id int64) (*model.GlobalVariables, error)
}
