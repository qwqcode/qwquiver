package exd

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

// GetAllExams 获取所有 Exam
func GetAllExams() []model.Exam {
	exams := []model.Exam{}
	lib.DB.Table(lib.ExamConfTb).Find(&exams)

	// 根据时间降序
	sort.Slice(exams, func(i, j int) bool {
		if exams[i].Date == "" || exams[j].Date == "" {
			return false
		}
		t1, err := dateparse.ParseAny(exams[i].Date)
		t2, err := dateparse.ParseAny(exams[j].Date)
		if err != nil {
			return false
		}
		return t1.After(t2)
	})

	return exams
}

// GetExam 获取 Exam
func GetExam(name string) *model.Exam {
	var exam model.Exam
	lib.DB.Table(lib.ExamConfTb).First(&exam, "name = ?", name)
	return &exam
}

// GetExamName 获取 Exam 名（非 tableName）
func GetExamName(name string) string {
	return strings.TrimPrefix(name, lib.ExamTbPf)
}

// GetExamTableName 获取 Exam 数据表名
func GetExamTableName(name string) string {
	return lib.ExamTbPf + strings.TrimPrefix(name, lib.ExamTbPf)
}

// HasExam 判断 Exam 是否存在
func HasExam(name string) bool {
	return lib.DB.Migrator().HasTable(GetExamTableName(name))
}

// RemoveExam 删除 Exam
func RemoveExam(name string) error {
	exam := GetExam(name)
	if exam == nil {
		return errors.New("未找到 Exam: " + name)
	}
	lib.DB.Table(lib.ExamConfTb).Delete(exam)
	return lib.DropTable(GetExamTableName(name))
}

// FilterScores 查询指定的成绩数据
func FilterScores(query *gorm.DB, rawMatchCond map[string]string, regMode bool) *gorm.DB {
	if regMode {
		i := 0
		for key, val := range rawMatchCond {
			if i == 0 {
				query = query.Where(key+` LIKE ?`, `%`+val+`%`)
			} else {
				query = query.Or(key+` LIKE ?`, `%`+val+`%`)
			}
			i++
		}
	} else {
		query = query.Where(rawMatchCond)
	}

	return query
}

// FilterScoresByRegStr 模糊查询指定的成绩数据
func FilterScoresByRegStr(query *gorm.DB, regStr string) *gorm.DB {
	return FilterScores(query, map[string]string{
		"NAME":   regStr,
		"SCHOOL": regStr,
		"CLASS":  regStr,
		"CODE":   regStr,
	}, true)
}

// GetAllExamGrps 获取所有 Exam 的 Grp
func GetAllExamGrps() []string {
	grps := []string{}
	for _, exam := range GetAllExams() {
		if exam.Grp != "" {
			if !funk.ContainsString(grps, exam.Grp) {
				grps = append(grps, exam.Grp) // 追加不重复的 Grp
			}
		}
	}
	return grps
}

// GetExamsByGrp 获取指定 Grp 的所有 Exam
func GetExamsByGrp(grp string) []model.Exam {
	exams := []model.Exam{}
	for _, exam := range GetAllExams() {
		if exam.Grp == grp {
			exams = append(exams, exam)
		}
	}
	return exams
}

// GetExamConfJSONStr 获取考试配置数据 JSON 字符串，beauty=是否格式化JSON
func GetExamConfJSONStr(name string, beauty bool) string {
	examConf := GetExam(name)
	if examConf == nil {
		return ""
	}
	if beauty {
		json, _ := json.MarshalIndent(examConf, "", "    ")
		return string(json)
	}

	json, _ := json.Marshal(examConf)
	return string(json)
}

// UpdateExamConf 修改考试配置
func UpdateExamConf(examConf *model.Exam) error {
	return lib.DB.Table(lib.ExamConfTb).Save(examConf).Error
}

// UpdateExamConfByJSON 修改考试配置，通过 JSON 数据
func UpdateExamConfByJSON(examConf *model.Exam, jsonStr string) error {
	var nExamConf model.Exam
	// testData := []byte(`{"Name":"测试","Grp":"233","Label":"233","Subj":["DL"],"SubjFullScore":{"DL":100},"Date":"dd","Note":"ww"}`)
	if err := json.Unmarshal([]byte(jsonStr), &nExamConf); err != nil {
		panic(err)
	}

	examConf.Grp = nExamConf.Grp
	examConf.Label = nExamConf.Label
	examConf.Subj = nExamConf.Subj
	examConf.SubjFullScore = nExamConf.SubjFullScore
	examConf.Date = nExamConf.Date
	examConf.Note = nExamConf.Note
	// TODO ... 比较 dirty 的代码，不够 flexible，先将就这样

	return UpdateExamConf(examConf)
}
