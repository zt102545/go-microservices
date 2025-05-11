package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type PostgresConfig struct {
	Addr        map[string]string `json:"addr,optional"`
	MaxIdleConn int               `json:"max_idle_conn,default=100"`
	MaxOpenConn int               `json:"max_open_conn,default=200"`
	MaxIdleTime int               `json:"max_idle_time,default=30"`
	db          map[string]*sql.DB
}

func (d PostgresConfig) Init() map[string]*sql.DB {

	d.db = make(map[string]*sql.DB)
	for k, addr := range d.Addr {

		db, err := sql.Open("postgres", addr)
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
		d.db[k] = db
	}
	return d.db
}

func (d PostgresConfig) Close() {
	if d.db != nil {
		for _, db := range d.db {
			err := db.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
