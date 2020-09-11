package model

import (
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"gorm.io/gorm"
)

const _examTbPf string = "EXAM_SCORES_" // 用于存放 Exam 数据的数据表名前缀

// Exam 考试配置数据 模型
type Exam struct {
	Name          string `gorm:"primaryKey"` // Same as a Bucket's Name
	Grp           string // 分类
	Label         string // 实际显示的名称
	Subj          string // 考试包含科目
	SubjFullScore string // 每科的满分分数
	Date          string // 考试日期 eg. "2020-08-07"
	Note          string // 备注
}

// GetName 获取 Exam 名
func (exam Exam) GetName() string {
	return exam.Name
}

// GetTableName 获取 Exam 数据表名
func (exam Exam) GetTableName() string {
	return lib.ExamTbPf + exam.Name
}

// NewQuery 新建 Exam 查询
func (exam Exam) NewQuery() *gorm.DB {
	return lib.DB.Table(lib.ExamTbPf + exam.Name)
}

// GetTable 获取存放成绩的数据表
func (exam Exam) GetTable() *gorm.DB {
	return exam.NewQuery()
}

// CountScores 获取成绩数据条数（考生数量）
func (exam Exam) CountScores() (num int64) {
	lib.DB.Table(lib.ExamTbPf + exam.Name).Count(&num)
	return
}

// GetSubjects 获取所有考试科目
func (exam Exam) GetSubjects() (subjects []string, err error) {
	if exam.Subj != "" {
		if err := utils.JSONDecode(exam.Subj, &subjects); err != nil {
			return []string{}, err
		}
	}
	return
}
