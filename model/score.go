package model

// Score 成绩模型
type Score struct {
	ID     int     `storm:"id,increment"` // 编号
	Name   string  `storm:"index"`        // 姓名
	Code   string  `storm:"index"`        // 考号
	School string  `storm:"index"`        // 学校
	Class  string  `storm:"index"`        // 班级
	Total  float64 // 总分
	Rank   int     // 排名

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

	ZKRank int // 主排
	LKRank int // 理排
	WKRank int // 文排

	LZRank int // 理综排
	WZRank int // 文综排
}

// ScoreFieldTransMap 字段名 => 中文名
var ScoreFieldTransMap map[string]string = map[string]string{
	"ID": "编号", "Name": "姓名", "Code": "考号", "School": "学校",
	"Class": "班级", "Total": "总分", "Rank": "排名",

	"YW": "语文", "SX": "数学", "YY": "英语",
	"WL": "物理", "HX": "化学", "SW": "生物",
	"ZZ": "政治", "LS": "历史", "DL": "地理",

	"LZ": "理综", "WZ": "文综",
	"ZK": "主科", "LK": "理科", "WK": "文科",

	"ZKRank": "主排", "LKRank": "理排", "WKRank": "文排",
	"LZRank": "理综排", "WZRank": "文综排",
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
var SFieldExtRank []string = []string{"ZKRank", "LKRank", "WKRank", "LZRank", "WZRank"}

// SFieldExtSum 拓展求和字段
var SFieldExtSum []string = []string{"ZK", "LZ", "WZ", "LK", "WK"}
