package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DbConfig struct {
	Addr        string `json:"addr,optional"`
	MaxIdleConn int    `json:"max_idle_conn,default=100"`
	MaxOpenConn int    `json:"max_open_conn,default=200"`
	MaxIdleTime int    `json:"max_idle_time,default=30"`
}

var db *sql.DB

func (d DbConfig) Init() *sql.DB {
	var err error
	db, err = sql.Open("mysql", d.Addr)
	if err != nil {
		log.Fatalf("initDB failed: %s", err.Error())
		return nil
	}

	db.SetMaxOpenConns(d.MaxOpenConn)
	db.SetMaxIdleConns(d.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(d.MaxIdleTime) * time.Second)
	err = db.Ping()
	if err != nil {
		log.Fatalf("initDB failed: %s", err.Error())
		return nil
	}
	return db
}

func (d DbConfig) Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}
