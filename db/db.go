package db

import (
	"database/sql"
	"sync"
)

const (
	path = "./db.sqlite3"
)

type PanDb struct {
	DB *sql.DB
	mu sync.Mutex
}

var instance *PanDb

func GetCacheInstance() *PanDb {
	if instance == nil {
		instance = &PanDb{
			DB: new(sql.DB),
		}
	}
	instance.DB = openDB()
	initDB(instance.DB)
	return instance
}

// openDB 打开数据库
func openDB() *sql.DB {
	//打开数据库，如果不存在，则创建
	db, err := sql.Open("sqlite3", path)
	checkErr(db, err)
	return db
}

// initDB 初始化数据库
func initDB(db *sql.DB) {
	//创建表
	sqlTable := `
	CREATE TABLE IF NOT EXISTS fileinfo(
		fs_id INTEGER PRIMARY KEY,
		p_id INTEGER NULL,
		app_id INTEGER,
		uk INTEGER,
		msg_id INTEGER,
		group_id INTEGER,
		category INTEGER,
		path TEXT,
		isdir INTEGER,
		server_filename TEXT,
		server_ctime INTEGER,
		server_mtime INTEGER,
		size INTEGER,
		dlink TEXT
	);
			`
	db.Exec(sqlTable)
}
