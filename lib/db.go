package lib

import (
	"errors"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/qwqcode/qwquiver/config"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
	"go.etcd.io/bbolt"
)

// DB is database
var DB *storm.DB

// 用于存放成绩数据的 Bucket 名称前缀
const scoreBucketPrefix string = "QwScore_"

// OpenDb 打开数据库
func OpenDb(dbFile string) (err error) {
	DB, err = storm.Open(config.Instance.DbFile)
	return
}

// CloseDb 关闭数据
func CloseDb() (err error) {
	return DB.Close()
}

// GetScoreBucketNames 获取所有 ScoreBucket 的 Name
func GetScoreBucketNames() (names []string) {
	_ = DB.Bolt.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			bn := string(name)
			if strings.HasPrefix(bn, scoreBucketPrefix) { // 以前缀开头 score bucket
				bn = strings.TrimPrefix(bn, scoreBucketPrefix) // 去掉前缀
				names = append(names, bn)
			}
			return nil
		})
	})
	return
}

// GetIsScoreBucketExist 判断 ScoreBucket 是否存在
func GetIsScoreBucketExist(name string) bool {
	name = strings.TrimPrefix(name, scoreBucketPrefix) // 若前缀存在则去掉前缀
	arr := GetScoreBucketNames()
	return funk.ContainsString(arr, name)
}

// CreateScoreBucket 创建新的 ScoreBucket
func CreateScoreBucket(name string) (err error) {
	if GetIsScoreBucketExist(name) {
		err = errors.New("创建 ScoreBucket 失败，名为 '" + name + "' 的 Bucket 已存在，不能重复创建")
		return
	}
	err = DB.From(scoreBucketPrefix + name).Init(&model.Score{})
	return
}

// GetScoreBucket 获取 ScoreBucket
func GetScoreBucket(name string) storm.Node {
	name = strings.TrimPrefix(name, scoreBucketPrefix)
	return DB.From(scoreBucketPrefix + name)
}

// FilterScores 查询指定的成绩数据
func FilterScores(bucketName string, matchCond map[string]interface{}, regMode bool) storm.Query {
	bucket := GetScoreBucket(bucketName)
	matchers := []q.Matcher{}

	for key, val := range matchCond {
		if regMode {
			matchers = append(matchers, q.Re(key, val.(string)))
		} else {
			matchers = append(matchers, q.Eq(key, val))
		}
	}

	if regMode {
		return bucket.Select(q.Or(matchers...))
	}
	return bucket.Select(matchers...)
}

// FilterScoresRegStr 查询指定的成绩数据
func FilterScoresRegStr(bucketName string, regStr string) storm.Query {
	return FilterScores(bucketName, map[string]interface{}{
		"Name":   regStr,
		"School": regStr,
		"Class":  regStr,
		"Code":   regStr,
	}, true)
}
