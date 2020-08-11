package lib

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cheggaaa/pb"
	"github.com/oleiade/reflections"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// ImportExcel 数据导入
func ImportExcel(examName string, filename string, examConfJSON string) {
	examName = strings.TrimSpace(examName)
	filename = strings.TrimSpace(filename)
	var examConf *model.ExamConf

	if filename == "" {
		logrus.Error("ExcelImporter: Excel 文件路径不能为空")
		return
	}

	if fileExt := filepath.Ext(filename); fileExt != ".xlsx" {
		logrus.Error("ExcelImporter: 仅支持导入 .xlsx 文件，'" + fileExt + "' 格式不被支持")
		return
	}

	if examName == "" { // 若 examName 为空，默认使用 Excel 文件名
		fileBasename := filepath.Base(filename)
		examName = strings.TrimSuffix(fileBasename, filepath.Ext(fileBasename))
	}

	// 查询是否相同 examName 的 bucket 已存在
	if IsExamExist(examName) {
		logrus.Error("ExcelImporter: 名称为 '" + examName + "' 的数据表已存在，无法重复导入；您可以执行 `qwquiver exam` 对现有 Exam 进行删除操作")
		return
	}

	// examConf 处理
	var jExamConf model.ExamConf
	if examConfJSON != "" {
		if err := json.Unmarshal([]byte(examConfJSON), &jExamConf); err != nil {
			logrus.Error("ExcelImporter: 解析 examConf JSON 发生错误 ", err)
			return
		}
		examConf = &jExamConf
	} else {
		examConf = &model.ExamConf{}
	}
	examConf.Name = examName

	logrus.Info("数据导入任务")
	logrus.Info("ExamName='" + examName + "', ExcelFile='" + filename + "'")

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

	optionalFields := utils.GetStructFields(&model.Score{})
	optionalFieldsTr := (funk.Values(model.ScoreFieldTransMap)).([]string) // 中文字段名

	logrus.Debugln("字段KEY ", optionalFields)
	logrus.Debugln("字段名 ", optionalFieldsTr)

	fmt.Println()
	fmt.Println(" - 表头 =", rows[0])

	// 读取表头
	fieldPos := map[string]int{} // 字段名 => Position
	fieldList := []string{}      // 所有字段
	for pos, val := range rows[0] {
		if funk.ContainsString(optionalFields, val) {
			fieldPos[val] = pos
			fieldList = append(fieldList, val)
		} else if funk.ContainsString(optionalFieldsTr, val) {
			// 若表头为中文名，则先翻译为原始字段名
			val, _ := funk.FindKey(model.ScoreFieldTransMap, func(trN string) bool {
				return trN == val
			})
			if val != "" {
				fieldPos[val.(string)] = pos
				fieldList = append(fieldList, val.(string))
			}
		}
	}
	subjList := funk.IntersectString(fieldList, model.SFieldSubj) // 考试科目
	examConf.Subj = subjList

	fmt.Println(" - 表头 POSITION =", fieldPos)
	fmt.Println(" - 考试科目 =", subjList)
	fmt.Println()

	scList := &[](*model.Score){}
	schoolList := map[string][]string{} // 学校&班级列表

	// 读取数据
	rankAbleFn := []string{} // 可排序的字段 字段名
	for _, row := range rows[1:] {
		sc := &model.Score{}

		for name, pos := range fieldPos {
			// 将表格值添加到 sc 中
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

		// 记录学校班级
		if sc.SCHOOL != "" {
			if schoolList[sc.SCHOOL] == nil {
				schoolList[sc.SCHOOL] = []string{}
			}
			if sc.CLASS != "" && !funk.Contains(schoolList[sc.SCHOOL], sc.CLASS) {
				schoolList[sc.SCHOOL] = append(schoolList[sc.SCHOOL], sc.CLASS)
			}
		}

		// 求和
		sc.ZK = sc.YW + sc.SX + sc.YY    // 主科
		sc.LZ = sc.WL + sc.HX + sc.SW    // 理综
		sc.WZ = sc.ZZ + sc.LS + sc.DL    // 文综
		sc.LK = sc.ZK + sc.LZ            // 理科
		sc.WK = sc.ZK + sc.WZ            // 文科
		sc.TOTAL = sc.ZK + sc.LZ + sc.WZ // 总分

		// 收集可排序的字段
		for _, field := range model.SFieldRankAble {
			fieldVal, err := reflections.GetField(sc, field)
			if err == nil && fieldVal.(float64) > 0 { // 字段值存在
				if !funk.Contains(rankAbleFn, field) {
					rankAbleFn = append(rankAbleFn, field)
				}
			}
		}

		*scList = append(*scList, sc)
	}
	logrus.Info("总分数据已生成")
	fmt.Println()
	fmt.Println(" - 排名执行字段 =", rankAbleFn)
	fmt.Println(" 开始生成排名数据...")
	fmt.Println()

	// 生成排名数据 func
	Rank := func(scList []*model.Score, rankByF string, outputF string) {
		nScList := make([]*model.Score, len(scList))
		copy(nScList, scList)

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
				rawSc := funk.Find(scList, func(x *model.Score) bool {
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
	Rank(*scList, "TOTAL", "RANK")
	for _, rankByF := range rankAbleFn {
		Rank(*scList, rankByF, rankByF+"_RANK") // 相对于全部数据的排名
	}

	for school, classes := range schoolList {
		// 相对于学校的排名
		scListRtSchool := funk.Filter(*scList, func(sc *model.Score) bool {
			return sc.SCHOOL == school
		}).([]*model.Score)
		Rank(scListRtSchool, "TOTAL", "SCHOOL_RANK")
		for _, rankByF := range rankAbleFn {
			Rank(scListRtSchool, rankByF, rankByF+"_SCHOOL_RANK")
		}

		// 相对于班级的排名
		for _, class := range classes {
			scListRtClass := funk.Filter(*scList, func(sc *model.Score) bool {
				return sc.SCHOOL == school && sc.CLASS == class
			}).([]*model.Score)
			Rank(scListRtClass, "TOTAL", "CLASS_RANK")
			for _, rankByF := range rankAbleFn {
				Rank(scListRtClass, rankByF, rankByF+"_CLASS_RANK")
			}
		}
	}

	logrus.Info("排名数据已生成")

	// 成绩从高到低排序
	sort.Slice(*scList, func(i, j int) bool {
		return (*scList)[i].TOTAL > (*scList)[j].TOTAL
	})

	fmt.Println()

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

	// 处理附加数据

	// 尝试获取每科的最高分数
	TryGetFullScore := func() (subjFullScore map[string]float64) {
		subjFullScore = map[string]float64{}
		RecordOnce := func(subj string, score float64) {
			if subj == "" || score <= 0 {
				return
			}

			var predScore float64 = 0
			if score > 110 {
				predScore = 150
			} else if score > 100 {
				predScore = 110
			} else if score > 90 {
				predScore = 100
			} else if score > 80 {
				predScore = 90
			} else if score > 50 {
				predScore = 60
			} else if score > 40 {
				predScore = 50
			}

			if predScore > subjFullScore[subj] {
				subjFullScore[subj] = predScore
			}
		}

		for _, sc := range *scList {
			for _, subj := range model.SFieldSubj {
				score := reflect.ValueOf(sc).Elem().FieldByName(subj).Float()
				RecordOnce(subj, score)
			}
		}

		return
	}
	subjFullScore := TryGetFullScore()
	if examConf.SubjFullScore == nil {
		examConf.SubjFullScore = map[string]float64{}
	}
	for subj, fullScore := range subjFullScore {
		if fullScore != 0 && examConf.SubjFullScore[subj] == 0 {
			examConf.SubjFullScore[subj] = fullScore
		}
	}

	logrus.Info("学科最高分数已获取")

	// 将数据导入数据库
	if err := CreateExam(examName); err != nil {
		logrus.Error("ExcelImporter: 创建 ScoreBucket 发生错误 ", err)
		return
	}

	// 保存 Exam 配置
	if err := SaveExamConf(examConf); err != nil {
		logrus.Error("ExcelImporter: 写入 ExamConf 发生错误 ", err)
		return
	}
	logrus.Info("ExamConf 已成功写入")

	// 准备开始一个 transaction
	bucket := GetExam(examName)
	tx, err := bucket.Begin(true) // https://github.com/asdine/storm#transactions
	if err != nil {
		logrus.Error("ExcelImporter: 准备 transaction 操作时发生错误 ", err)
	}
	defer tx.Rollback()

	fmt.Print("\n")
	saveBar := pb.StartNew(len(*scList))

	// 开始遍历导入数据库
	saveErr := []error{}
	itemCount := 0
	for _, sc := range *scList {
		err := tx.Save(sc)
		if err != nil {
			saveErr = append(saveErr, err)
		}
		saveBar.Add(1)
		itemCount++
	}
	tx.Commit()

	saveBar.Finish()
	fmt.Print("\n\n")
	logrus.Info("成功导入 ", itemCount, " 条数据")

	// 错误处理
	if len(saveErr) > 0 {
		logrus.Error("ExcelImporter: 保存 Score 过程中发生错误")
		for _, err := range saveErr {
			logrus.Error("  ", err)
		}
		return
	}

	fmt.Print("\n\n")
	fmt.Println("ExamConf = '" + GetExamConfJSONStr(examName, false) + "'")
	fmt.Print("\n\n")
	fmt.Println("您可以执行以下命令，对 Exam 进行修改：")
	fmt.Println("  - 更改配置：qwquiver exam config set \"" + examName + "\" -h")
	fmt.Println("  - 获取配置：qwquiver exam config get \"" + examName + "\" -i")
	fmt.Println("  - 执行删除：qwquiver exam remove \"" + examName + "\" --force")

	fmt.Println()
	logrus.Info("数据导入任务执行完毕")

	// 尝试读取数据库数据
	//var queryScores []model.Score
	//lib.GetScoreBucket(examName).All(&queryScores)
	// fmt.Println(len(queryScores))
}
