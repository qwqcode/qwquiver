package model

// ExamConf 考试配置数据 模型
type ExamConf struct {
	Name          string `gorm:"primaryKey"` // Same as a Bucket's Name
	Grp           string // 分类
	Label         string // 实际显示的名称
	Subj          string // 考试包含科目
	SubjFullScore string // 每科的满分分数
	Date          string // 考试日期 eg. "2020-08-07"
	Note          string // 备注
}
