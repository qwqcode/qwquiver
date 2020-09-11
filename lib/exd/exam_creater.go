package exd

import (
	"errors"

	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/model"
)

// CreateExam 创建新的 Exam
func CreateExam(exam *model.Exam) (err error) {
	if exam.Name == "" {
		return errors.New("Exam 的 Name 不能为空")
	}

	if HasExam(exam.Name) {
		err = errors.New("创建 Exam 失败，名为 '" + exam.Name + "' 的 Exam 已存在，不能重复创建")
		return
	}

	_SaveExamConf(exam)

	tableName := GetExamTableName(exam.Name)
	model := &model.Score{}
	if err := lib.DB.Table(tableName).Migrator().CreateTable(model); err != nil {
		panic(err)
	}
	// TODO: tableName非法字符的处理
	lib.DB.Exec("CREATE INDEX `idx_name_" + tableName + "` ON `" + tableName + "` (name)")
	lib.DB.Exec("CREATE INDEX `idx_code_" + tableName + "` ON `" + tableName + "` (code)")
	lib.DB.Exec("CREATE INDEX `idx_school_" + tableName + "` ON `" + tableName + "` (school)")
	lib.DB.Exec("CREATE INDEX `idx_class_" + tableName + "` ON `" + tableName + "` (class)")
	return
}

// SaveExamConf 保存新的考试配置
func _SaveExamConf(examConf *model.Exam) error {
	if examConf.Name == "" {
		return errors.New("examConf 的 Name 不能为空")
	}

	if !lib.DB.Migrator().HasTable(lib.ExamConfTb) {
		if err := lib.DB.Migrator().CreateTable(&model.Exam{}); err != nil {
			panic(err)
		}
		if err := lib.DB.Migrator().RenameTable(&model.Exam{}, lib.ExamConfTb); err != nil {
			panic(err)
		}
	}

	var findExamConf model.Exam
	_ = lib.DB.Table(lib.ExamConfTb).Find(&findExamConf, "Name = ?", examConf.Name)
	if (findExamConf != model.Exam{}) {
		// 若已存在，则删除旧数据
		if r := lib.DB.Table(lib.ExamConfTb).Delete(findExamConf); r.Error != nil {
			panic(r.Error)
		}
	}

	return lib.DB.Table(lib.ExamConfTb).Create(examConf).Error
}
