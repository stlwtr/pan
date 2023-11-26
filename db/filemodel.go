package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// FileModel 文件模型
type FileModel struct {
	db              *sql.DB
	fs_id           int64
	p_id            int64
	app_id          int64
	uk              int64
	msg_id          int64
	group_id        int64
	category        int64
	path            string
	isdir           int64
	server_filename string
	server_ctime    int64
	server_mtime    int64
	size            int64
	dlink           string
}

func NewFileModel() *FileModel {
	return &FileModel{}
}

func checkErr(data interface{}, err error) (interface{}, error) {
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return data, err
}

// insert 新增
func (fm FileModel) insert() (sql.Result, error) {
	GetCacheInstance().mu.Lock()
	defer GetCacheInstance().mu.Unlock()
	db := GetCacheInstance().DB

	stmt, err := db.Prepare("insert into fileinfo(fs_id, p_id, app_id, uk, msg_id, group_id, category, path, isdir, server_filename, server_ctime, size, dlink) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	checkErr(stmt, err)
	res, err := stmt.Exec(fm.fs_id, fm.p_id, fm.app_id, fm.uk, fm.msg_id, fm.group_id, fm.category, fm.path, fm.isdir, fm.server_filename, fm.server_ctime, fm.server_mtime, fm.size, fm.dlink)
	checkErr(res, err)
	return res, nil
}

// update	更新文件
func (fm FileModel) update(fs_id int64) int64 {
	GetCacheInstance().mu.Lock()
	defer GetCacheInstance().mu.Unlock()
	db := GetCacheInstance().DB

	stmt, err := db.Prepare("update fileinfo set p_id=?, app_id=?, uk=?, msg_id=?, group_id=?, category=?, path=?, isdir=?, server_filename=?, server_ctime=?, server_mtime=?, size=?, dlink=? where fs_id=?")
	checkErr(stmt, err)
	res, err := stmt.Exec(fm.p_id, fm.app_id, fm.uk, fm.msg_id, fm.group_id, fm.category, fm.path, fm.isdir, fm.server_filename, fm.server_ctime, fm.server_mtime, fm.size, fm.dlink, fm.fs_id)
	checkErr(res, err)
	affect, err := res.RowsAffected()
	checkErr(affect, err)
	return affect
}

// query 查询
func (fm FileModel) query(condition string) ([]FileModel, error) {
	GetCacheInstance().mu.Lock()
	defer GetCacheInstance().mu.Unlock()
	db := GetCacheInstance().DB

	q := "select * from fileinfo"
	if len(condition) > 0 {
		q = q + " where " + condition
	}
	rows, err := db.Query(q)
	checkErr(rows, err)
	var fileList = []FileModel{}
	for rows.Next() {
		var fm = FileModel{}
		err = rows.Scan(&fm.fs_id, &fm.p_id, &fm.app_id, &fm.uk, &fm.msg_id, &fm.group_id, &fm.category, &fm.path, &fm.isdir, &fm.server_filename, &fm.server_ctime, &fm.server_mtime, &fm.size, &fm.dlink)
		checkErr(nil, err)
		fileList = append(fileList, fm)
	}
	rows.Close()
	return fileList, nil
}
