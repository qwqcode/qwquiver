package model

// Score 成绩模型
type Score struct {
	ID     int    `storm:"id,increment"` // 编号
	Name   string `storm:"index"`        // 姓名
	Code   string `storm:"index"`        // 考号
	School string `storm:"index"`        // 学校
	Class  string `storm:"index"`        // 班级
	Type   int    `storm:"index"`        // 类型
	Total  int    // 总分
	Rank   int    // 排名

	YW int // 语文
	SX int // 数学
	YY int // 英语

	WL int // 物理
	HX int // 化学
	SW int // 生物

	ZZ int // 政治
	LS int // 历史
	DL int // 地理

	LZ int // 理综 (物+化+生)
	WZ int // 文综 (政+历+地)

	ZK int // 主科 (语+数+英)
	LK int // 理科 (语数英+理综)
	WK int // 文科 (语数英+理综)

	ZKRank int // 主排
	LKRank int // 理排
	WKRank int // 文排

	LZRank int // 理综排
	WZRank int // 文综排
}
