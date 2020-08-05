package main

import (
	"github.com/qwqcode/qwquiver/cli/cmd"
)

func init() {

	// // 打开数据库
	// if err := lib.OpenDb(config.Instance.DbFile); err != nil {
	// 	logrus.Error(err)
	// }
}

func main() {
	cmd.Execute()
}
