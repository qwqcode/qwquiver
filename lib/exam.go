package lib

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
)

// 用于存放考试成绩数据的 Bucket 名称前缀
const examBucketPrefix string = "EXAM_SCORES_"
const examConfBucketName string = "EXAM_CONF"

// GetAllExamNames 获取数据库中每场考试的名称
func GetAllExamNames() (allNames []string) {
	for _, bName := range GetAllBucketNames() {
		if strings.HasPrefix(bName, examBucketPrefix) { // 以某前缀开头 bucket
			examName := strings.TrimPrefix(string(bName), examBucketPrefix) // 去掉前缀
			allNames = append(allNames, examName)
		}
	}
	return
}

// IsExamExist 判断 Exam 是否存在
func IsExamExist(name string) bool {
	name = strings.TrimPrefix(name, examBucketPrefix) // 若前缀存在则去掉前缀
	return funk.ContainsString(GetAllExamNames(), name)
}

// CreateExam 创建新的 Exam
func CreateExam(name string) (err error) {
	if IsExamExist(name) {
		err = errors.New("创建 Exam 失败，名为 '" + name + "' 的 Exam 已存在，不能重复创建")
		return
	}
	err = DB.From(examBucketPrefix + name).Init(&model.Score{})
	return
}

// GetExam 获取 Exam
func GetExam(name string) storm.Node {
	name = strings.TrimPrefix(name, examBucketPrefix)
	return DB.From(examBucketPrefix + name)
}

// RemoveExam 删除 Exam
func RemoveExam(name string) error {
	name = strings.TrimPrefix(name, examBucketPrefix)
	return RemoveBucket(examBucketPrefix + name)
}

// FilterScores 查询指定的成绩数据
func FilterScores(exam storm.Node, matchCond map[string]interface{}, regMode bool) storm.Query {
	matchers := []q.Matcher{}

	for key, val := range matchCond {
		if regMode {
			matchers = append(matchers, q.Re(key, val.(string)))
		} else {
			matchers = append(matchers, q.Eq(key, val))
		}
	}

	if regMode {
		return exam.Select(q.Or(matchers...))
	}
	return exam.Select(matchers...)
}

// FilterScoresByRegStr 模糊查询指定的成绩数据
func FilterScoresByRegStr(exam storm.Node, regStr string) storm.Query {
	return FilterScores(exam, map[string]interface{}{
		"Name":   regStr,
		"School": regStr,
		"Class":  regStr,
		"Code":   regStr,
	}, true)
}

// GetExamScoreLen 获取考试成绩数量（考生数量）
func GetExamScoreLen(name string) (num int) {
	num, _ = GetExam(name).Count(&model.Score{})
	return
}

// GetExamConf 获取考试配置数据
func GetExamConf(name string) *model.ExamConf {
	var conf model.ExamConf
	if err := DB.From(examConfBucketName).One("Name", name, &conf); err != nil {
		return nil
	}
	return &conf
}

// GetExamConfJSONStr 获取考试配置数据 JSON 字符串，beauty=是否格式化JSON
func GetExamConfJSONStr(name string, beauty bool) string {
	examConf := GetExamConf(name)
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

// SaveExamConf 保存新的考试配置
func SaveExamConf(examConf *model.ExamConf) error {
	if examConf.Name == "" {
		return errors.New("examConf 的 Name 不能为空")
	}

	var findExamConf *model.ExamConf
	_ = DB.From(examConfBucketName).One("Name", examConf.Name, &findExamConf)
	if findExamConf != nil {
		// 若已存在，则删除旧数据
		err := DB.From(examConfBucketName).DeleteStruct(findExamConf)
		if err != nil {
			return err
		}
	}

	return DB.From(examConfBucketName).Save(examConf)
}

// UpdateExamConf 修改考试配置
func UpdateExamConf(examConf *model.ExamConf) error {
	return DB.From(examConfBucketName).Update(examConf)
}

// UpdateExamConfByJSON 修改考试配置，通过 JSON 数据
func UpdateExamConfByJSON(examConf *model.ExamConf, jsonStr string) error {
	var nExamConf model.ExamConf
	// testData := []byte(`{"Name":"测试","Grp":"233","Label":"233","SubjFullScore":{"DL":100},"Date":"dd","Note":"ww"}`)
	if err := json.Unmarshal([]byte(jsonStr), &nExamConf); err != nil {
		return err
	}

	examConf.Grp = nExamConf.Grp
	examConf.Label = nExamConf.Label
	examConf.SubjFullScore = nExamConf.SubjFullScore
	examConf.Date = nExamConf.Date
	examConf.Note = nExamConf.Note
	// TODO ... 比较 dirty 的代码，先将就

	return UpdateExamConf(examConf)
}
