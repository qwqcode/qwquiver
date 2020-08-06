package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

func ImportExcel(filename string) {
	if fileExt := filepath.Ext(filename); fileExt != ".xlsx" {
		logrus.Error("ExcelImporter: 仅支持导入 .xlsx 文件，'" + fileExt + "' 格式不被支持")
		return
	}

	file, err := excelize.OpenFile(filename)
	if err != nil {
		logrus.Error("ExcelImporter: Excel 打开失败 ", err)
		return
	}

	rows := file.GetRows(file.GetSheetName(1))
	if len(rows) < 1 {
		logrus.Error("ExcelImporter: Excel 数据不能为空")
		return
	}

	f := map[string]int{}
	optionalFields := utils.GetStructFields(&model.Score{})
	optionalFieldsTr := (funk.Values(model.ScoreFieldTransMap)).([]string) // 中文字段名

	fmt.Println("1", optionalFields)
	fmt.Println("2", optionalFieldsTr)

	// 读取表头
	for pos, name := range rows[0] {
		if funk.ContainsString(optionalFields, name) {
			f[name] = pos
		} else if funk.ContainsString(optionalFieldsTr, name) {
			name, _ := funk.FindKey(model.ScoreFieldTransMap, func(trN string) bool {
				return trN == name
			})
			f[name.(string)] = pos
		}
	}

	fmt.Println("3", f)

	// 读取数据
	for _, row := range rows[1:] {
		score := &model.Score{}

		for name, pos := range f {
			rf := reflect.ValueOf(score).Elem().FieldByName(name)
			switch reflect.ValueOf(model.Score{}).FieldByName(name).Kind() {
			case reflect.String:
				rf.SetString(row[pos])
			case reflect.Float64:
				val, _ := strconv.ParseFloat(row[pos], 64)
				rf.SetFloat(val)
			}
		}

		fmt.Println(score)
		os.Exit(0)
	}
}
