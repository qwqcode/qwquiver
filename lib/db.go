package lib

import (
	"github.com/asdine/storm"
	"github.com/qwqcode/qwquiver/config"
	"github.com/thoas/go-funk"
	"go.etcd.io/bbolt"
)

// DB is database
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

// GetAllBucketNames 获取所有 Bucket 的名称
func GetAllBucketNames() (allNames []string) {
	_ = DB.Bolt.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			allNames = append(allNames, string(name))
			return nil
		})
	})
	return
}

// IsBucketExits 判断 Bucket 是否存在
func IsBucketExits(name string) bool {
	return funk.ContainsString(GetAllBucketNames(), name)
}

// RemoveBucket 删除 Bucket
func RemoveBucket(name string) error {
	return DB.Drop(name)
}
