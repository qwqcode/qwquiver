package model

// Score 成绩模型
type Score struct {
	ID          int     `gorm:"primaryKey;autoIncrement"` // 编号
	NAME        string  ``                                // 姓名
	CODE        string  ``                                // 考号
	SCHOOL      string  ``                                // 学校
	CLASS       string  ``                                // 班级
	TOTAL       float64 // 总分
	RANK        int     // 排名
	SCHOOL_RANK int     // 校排名
	CLASS_RANK  int     // 校排名

	YW float64 // 语文
	SX float64 // 数学
	YY float64 // 英语

	WL float64 // 物理
	HX float64 // 化学
	SW float64 // 生物

	ZZ float64 // 政治
	LS float64 // 历史
	DL float64 // 地理

	LZ float64 // 理综 (物+化+生)
	WZ float64 // 文综 (政+历+地)

	ZK float64 // 主科 (语+数+英)
	LK float64 // 理科 (语数英+理综)
	WK float64 // 文科 (语数英+理综)

	ZK_RANK int // 主排
	LK_RANK int // 理排
	WK_RANK int // 文排
	LZ_RANK int // 理综排
	WZ_RANK int // 文综排
	YW_RANK int // 语文排
	SX_RANK int // 数学排
	YY_RANK int // 英语排
	WL_RANK int // 物理排
	HX_RANK int // 化学排
	SW_RANK int // 生物排
	ZZ_RANK int // 政治排
	LS_RANK int // 历史排
	DL_RANK int // 地理排

	ZK_SCHOOL_RANK int // 主班排
	LK_SCHOOL_RANK int // 理班排
	WK_SCHOOL_RANK int // 文班排
	LZ_SCHOOL_RANK int // 理综班排
	WZ_SCHOOL_RANK int // 文综班排
	YW_SCHOOL_RANK int // 语文班排
	SX_SCHOOL_RANK int // 数学班排
	YY_SCHOOL_RANK int // 英语班排
	WL_SCHOOL_RANK int // 物理班排
	HX_SCHOOL_RANK int // 化学班排
	SW_SCHOOL_RANK int // 生物班排
	ZZ_SCHOOL_RANK int // 政治班排
	LS_SCHOOL_RANK int // 历史班排
	DL_SCHOOL_RANK int // 地理班排

	ZK_CLASS_RANK int // 主班排
	LK_CLASS_RANK int // 理班排
	WK_CLASS_RANK int // 文班排
	LZ_CLASS_RANK int // 理综班排
	WZ_CLASS_RANK int // 文综班排
	YW_CLASS_RANK int // 语文班排
	SX_CLASS_RANK int // 数学班排
	YY_CLASS_RANK int // 英语班排
	WL_CLASS_RANK int // 物理班排
	HX_CLASS_RANK int // 化学班排
	SW_CLASS_RANK int // 生物班排
	ZZ_CLASS_RANK int // 政治班排
	LS_CLASS_RANK int // 历史班排
	DL_CLASS_RANK int // 地理班排
}

// ScoreFieldTransMap 字段名 => 中文名
var ScoreFieldTransMap map[string]string = map[string]string{
	"ID": "编号", "NAME": "姓名", "CODE": "考号", "SCHOOL": "学校",
	"CLASS": "班级", "TOTAL": "总分", "RANK": "排名",

	"YW": "语文", "SX": "数学", "YY": "英语",
	"WL": "物理", "HX": "化学", "SW": "生物",
	"ZZ": "政治", "LS": "历史", "DL": "地理",

	"LZ": "理综", "WZ": "文综",
	"ZK": "主科", "LK": "理科", "WK": "文科",

	"SCHOOL_RANK": "校排名", "CLASS_RANK": "班排名",
}

// SFieldSubj 所有学科字段名
var SFieldSubj []string = []string{"YW", "SX", "YY", "WL", "HX", "SW", "ZZ", "LS", "DL"}

// SFieldSubjZK 所有主要科目字段名
var SFieldSubjZK []string = []string{"YW", "SX", "YY"}

// SFieldSubjLK 所有理科字段名
var SFieldSubjLK []string = []string{"WL", "HX", "SW"}

// SFieldSubjWK 所有文科字段名
var SFieldSubjWK []string = []string{"ZZ", "LS", "DL"}

// SFieldExtRank 拓展排名字段
var SFieldExtRank []string = []string{"ZK_RANK", "LK_RANK", "WK_RANK", "LZ_RANK", "WZ_RANK"}

// SFieldRankAble 可进行排名的字段
var SFieldRankAble []string = []string{"ZK", "LK", "WK", "LZ", "WZ", "YW", "SX", "YY", "WL", "HX", "SW", "ZZ", "LS", "DL"}

// SFieldExtSum 拓展求和字段
var SFieldExtSum []string = []string{"ZK", "LZ", "WZ", "LK", "WK"}
