package lib

import (
	"github.com/qwqcode/qwquiver/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is database
var DB *gorm.DB

// OpenDb 打开数据库
func OpenDb(dbFile string) (err error) {
	DB, err = gorm.Open(sqlite.Open(config.Instance.DbFile), &gorm.Config{})
	return
}

// GetTables 获取所有数据表名称
func GetTables() (tables []string) {
	rows, err := DB.Raw("SELECT `name` FROM sqlite_master WHERE type='table';").Rows()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			panic(err)
		}

		tables = append(tables, table)
		// if table != "schema_migrations" {
		// 	tables = append(tables, table)
		// }
	}
	return
}

// HasTable 判断数据表是否存在
func HasTable(name string) bool {
	return DB.Migrator().HasTable(name)
}

// DropTable 删除数据表
func DropTable(name string) error {
	return DB.Migrator().DropTable(name)
}
