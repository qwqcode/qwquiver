package lib

import (
	"github.com/asdine/storm"
	"github.com/qwqcode/qwquiver/config"
)

var DB *storm.DB

// OpenDb 打开数据库
func OpenDb(dbFile string) (err error) {
	DB, err = storm.Open(config.Instance.DbFile)
	return
}

// CloseDb 关闭数据
func CloseDb() (err error) {
	return DB.Close()
}
