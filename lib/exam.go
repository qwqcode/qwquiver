package lib

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

const _examTbPf string = "EXAM_SCORES_" // 用于存放 Exam 数据的数据表名前缀
const _examConfTb string = "EXAM_CONF"  // ExamConf 数据表名

// GetAllExamNames 获取数据库中所有 Exam 的名称
func GetAllExamNames() (examNames []string) {
	for _, tbName := range GetTables() {
		if strings.HasPrefix(tbName, _examTbPf) { // 以某前缀开头 table
			name := strings.TrimPrefix(tbName, _examTbPf) // 去掉前缀
			examNames = append(examNames, name)
		}
	}
	return
}

// GetExamName 获取 Exam 名（非 tableName）
func GetExamName(name string) string {
	return strings.TrimPrefix(name, _examTbPf)
}

// GetExamTableName 获取 Exam 数据表名
func GetExamTableName(name string) string {
	return _examTbPf + strings.TrimPrefix(name, _examTbPf)
}

// HasExam 判断 Exam 是否存在
func HasExam(name string) bool {
	return DB.Migrator().HasTable(GetExamTableName(name))
}

// CreateExam 创建新的 Exam
func CreateExam(name string) (err error) {
	if HasExam(name) {
		err = errors.New("创建 Exam 失败，名为 '" + name + "' 的 Exam 已存在，不能重复创建")
		return
	}

	tableName := GetExamTableName(name)
	score := model.Score{}
	if err := DB.Migrator().CreateTable(&score); err != nil {
		panic(err)
	}
	if err := DB.Migrator().RenameTable(&score, tableName); err != nil {
		panic(err)
	}
	// TODO: index 搞不定
	// err = DB.Migrator().RenameIndex(&score, "idx_score_name", "idx_score_name_"+name)
	// err = DB.Migrator().RenameIndex(&score, "idx_score_code", "idx_score_code_"+name)
	// err = DB.Migrator().RenameIndex(&score, "idx_score_school", "idx_score_school_"+name)
	// err = DB.Migrator().RenameIndex(&score, "idx_score_class", "idx_score_class_"+name)
	if err != nil {
		panic(err)
	}
	return
}

// NewExamQuery 获取 Exam
func NewExamQuery(name string) *gorm.DB {
	return DB.Table(GetExamTableName(name))
}

// RemoveExam 删除 Exam
func RemoveExam(name string) error {
	return DropTable(GetExamTableName(name))
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

// GetExamScoreLen 获取 Exam 成绩数据数量（考生数量）
func GetExamScoreLen(name string) (num int64) {
	DB.Table(GetExamTableName(name)).Count(&num)
	return
}

// GetExamConf 获取 Exam 配置数据
func GetExamConf(name string) *model.ExamConf {
	var conf model.ExamConf
	DB.Table(_examConfTb).First(&conf, "name = ?", name)
	return &conf
}

// GetAllExamConf 获取所有 Exam 的 ExamConf
func GetAllExamConf() (examConfList []model.ExamConf) {
	examConfList = []model.ExamConf{}
	for _, name := range GetAllExamNames() {
		examConf := GetExamConf(name)
		if examConf != nil {
			examConfList = append(examConfList, *examConf)
		}
	}
	return
}

// GetAllExamGrps 获取所有已存在的 Grp
func GetAllExamGrps() (examGrpList []string) {
	examGrpList = []string{}
	for _, examName := range GetAllExamNames() {
		examConf := GetExamConf(examName)
		if examConf != nil && examConf.Grp != "" {
			if !funk.ContainsString(examGrpList, examConf.Grp) {
				examGrpList = append(examGrpList, examConf.Grp) // 追加不重复的 Grp
			}
		}
	}
	return
}

// ExamBean a bean for storing Exam's information
type ExamBean struct {
	Name string
	Conf *model.ExamConf
}

// GetExamsByGrp 获取指定 Grp 的所有 Exam
func GetExamsByGrp(grp string) (exams []ExamBean) {
	exams = []ExamBean{}
	for _, examName := range GetAllExamNames() {
		examConf := GetExamConf(examName)
		if examConf != nil && examConf.Grp == grp {
			exams = append(exams, ExamBean{
				Name: examName,
				Conf: examConf,
			})
		}
	}
	return
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

	if !DB.Migrator().HasTable(_examConfTb) {
		if err := DB.Migrator().CreateTable(&model.ExamConf{}); err != nil {
			panic(err)
		}
		if err := DB.Migrator().RenameTable(&model.ExamConf{}, _examConfTb); err != nil {
			panic(err)
		}
	}

	var findExamConf model.ExamConf
	_ = DB.Table(_examConfTb).Find(&findExamConf, "Name = ?", examConf.Name)
	if (findExamConf != model.ExamConf{}) {
		// 若已存在，则删除旧数据
		if r := DB.Table(_examConfTb).Delete(findExamConf); r.Error != nil {
			panic(r.Error)
		}
	}

	return DB.Table(_examConfTb).Create(examConf).Error
}

// UpdateExamConf 修改考试配置
func UpdateExamConf(examConf *model.ExamConf) error {
	return DB.Save(examConf).Error
}

// UpdateExamConfByJSON 修改考试配置，通过 JSON 数据
func UpdateExamConfByJSON(examConf *model.ExamConf, jsonStr string) error {
	var nExamConf model.ExamConf
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
