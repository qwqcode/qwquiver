package tools

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

func ImportExcel(dataName string, filename string) {
	dataName = strings.TrimSpace(dataName)
	filename = strings.TrimSpace(filename)

	if filename == "" {
		logrus.Error("ExcelImporter: Excel 文件路径不能为空")
		return
	}

	if fileExt := filepath.Ext(filename); fileExt != ".xlsx" {
		logrus.Error("ExcelImporter: 仅支持导入 .xlsx 文件，'" + fileExt + "' 格式不被支持")
		return
	}

	if dataName == "" { // 若 dataName 为空，默认使用 Excel 文件名
		fileBasename := filepath.Base(filename)
		dataName = strings.TrimSuffix(fileBasename, filepath.Ext(fileBasename))
	}

	// 查询是否相同 dataName 的 bucket 已存在
	if lib.GetIsScoreBucketExist(dataName) {
		logrus.Error("ExcelImporter: 名称为 '" + dataName + "' 的数据表已存在，无法重复导入；您可以加上 flag: '--force' 强制覆盖数据表，或使用 'qwquiver db' 对数据库进行管理")
		return
	}

	logrus.Info("ExcelImporter: 执行 Excel 数据导入任务，DataName='" + dataName + "', FileName='" + filename + "'")

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

	scList := &[](*model.Score){}

	// 读取数据
	rankAbleFn := []string{} // 可排序的字段 字段名
	for _, row := range rows[1:] {
		sc := &model.Score{}

		for name, pos := range f {
			rf := reflect.ValueOf(sc).Elem().FieldByName(name)
			switch reflect.ValueOf(model.Score{}).FieldByName(name).Kind() {
			case reflect.String:
				rf.SetString(row[pos])
			case reflect.Float64:
				val, _ := strconv.ParseFloat(row[pos], 64)
				rf.SetFloat(val)
			case reflect.Int:
				val, _ := strconv.Atoi(row[pos])
				rf.SetInt(int64(val))
			}
		}

		// 求和
		sc.ZK = sc.YW + sc.SX + sc.YY    // 主科
		sc.LZ = sc.WL + sc.HX + sc.SW    // 理综
		sc.WZ = sc.ZZ + sc.LS + sc.DL    // 文综
		sc.LK = sc.ZK + sc.LZ            // 理科
		sc.WK = sc.ZK + sc.WZ            // 文科
		sc.Total = sc.ZK + sc.LZ + sc.WZ // 总分

		// 收集可排序的字段
		{
			if sc.ZK > 0 && !funk.Contains(rankAbleFn, "ZK") {
				rankAbleFn = append(rankAbleFn, "ZK")
			}
			if sc.LZ > 0 && !funk.Contains(rankAbleFn, "LZ") {
				rankAbleFn = append(rankAbleFn, "LZ")
			}
			if sc.WZ > 0 && !funk.Contains(rankAbleFn, "WZ") {
				rankAbleFn = append(rankAbleFn, "WZ")
			}
			if sc.LK > 0 && !funk.Contains(rankAbleFn, "LK") {
				rankAbleFn = append(rankAbleFn, "LK")
			}
			if sc.WK > 0 && !funk.Contains(rankAbleFn, "WK") {
				rankAbleFn = append(rankAbleFn, "WK")
			}
		}

		*scList = append(*scList, sc)
	}

	// 生成排名数据 func
	rank := func(scList *[]*model.Score, rankByF string, outputF string) {
		nScList := make([]*model.Score, len(*scList))
		copy(nScList, *scList)

		// 成绩从大到小排序
		sort.Slice(nScList, func(i, j int) bool {
			iv := reflect.ValueOf(nScList[i]).Elem().FieldByName(rankByF).Float() // 当前元素
			jv := reflect.ValueOf(nScList[j]).Elem().FieldByName(rankByF).Float() // 下一个元素

			return iv > jv // 从大到小，降序；当前元素是否大于下一个元素，当前元素是否排前面 (true/false)
		})

		// 建立排名
		var tRank int = 1
		var tNum float64 = -1
		var tSameNum int = 1
		for _, sc := range nScList {
			num := reflect.ValueOf(sc).Elem().FieldByName(rankByF).Float()
			setRankData := func(rankVal int) {
				rawSc := funk.Find(*scList, func(x *model.Score) bool {
					return *x == *sc
				})
				reflect.ValueOf(rawSc).Elem().FieldByName(outputF).SetInt(int64(rankVal))
				// fmt.Println(rawSc.(*model.Score).Name, " ", rankVal)
			}

			if tNum == -1 { // 最高分初始化
				tNum = num
				setRankData(tRank)
				continue
			}

			if num == tNum {
				tSameNum++
			} else if num < tNum {
				tRank = tRank + tSameNum
				tNum = num
				tSameNum = 1
			}
			setRankData(tRank)
		}
	}

	// 执行排名
	rank(scList, "Total", "Rank")
	for _, rankByF := range rankAbleFn {
		rank(scList, rankByF, rankByF+"Rank")
	}

	// 成绩从高到低排序
	sort.Slice(*scList, func(i, j int) bool {
		return (*scList)[i].Total > (*scList)[j].Total
	})

	consoleOutput := func() {
		for _, sc := range *scList {
			s := ""
			for _, f := range utils.GetStructFields(&model.Score{}) {
				rf := reflect.ValueOf(sc).Elem().FieldByName(f)
				s += model.ScoreFieldTransMap[f] + ": "
				switch reflect.ValueOf(model.Score{}).FieldByName(f).Kind() {
				case reflect.String:
					s += rf.String()
				case reflect.Float64:
					s += fmt.Sprintf("%.6f", rf.Float())
				case reflect.Int:
					s += fmt.Sprintf("%d", rf.Int())
				}
				s += ", "
			}
			fmt.Println(s)
		}
	}
	if false {
		consoleOutput()
	}

	// 将数据导入数据库
	if err := lib.CreateScoreBucket(dataName); err != nil {
		logrus.Error("ExcelImporter: 创建 ScoreBucket 发生错误 ", err)
		return
	}
	bucket := lib.GetScoreBucket(dataName)
	saveErr := []error{}
	for _, sc := range *scList {
		err := bucket.Save(sc)
		if err != nil {
			saveErr = append(saveErr, err)
		}
	}

	if len(saveErr) > 0 {
		logrus.Error("ExcelImporter: 保存 Score 过程中发生错误")
		for _, err := range saveErr {
			logrus.Error("  ", err)
		}
		return
	}

	// 尝试读取数据库数据
	var queryScores []model.Score
	lib.GetScoreBucket(dataName).All(&queryScores)

	fmt.Println(len(queryScores))
}
