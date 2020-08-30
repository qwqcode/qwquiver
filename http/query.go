package http

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
	"gopkg.in/oleiade/reflections.v1"
)

func queryHandler(c echo.Context) error {
	examName := c.QueryParam("exam")
	whereJSONStr := c.QueryParam("where")
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")
	sortJSONStr := c.QueryParam("sort")

	var initConf Map
	isInitReq := c.QueryParam("init") != ""
	if isInitReq {
		// 若为初始化请求
		examMap := lib.GetAllExamConf()
		examGrpList := lib.GetAllExamGrps()
		fieldTransDict := model.ScoreFieldTransMap

		if len(examMap) == 0 {
			return RespError(c, "未找到任何考试数据，请导入数据")
		}

		if examName == "" {
			examName = lib.GetAllExamNames()[0] // 设置默认 exam
		}

		initConf = Map{
			"examMap":        examMap,
			"examGrpList":    examGrpList,
			"fieldTransDict": fieldTransDict,
		}
	}

	if !lib.HasExam(examName) {
		return RespError(c, "Exam 不存在")
	}
	query := lib.NewExamQuery(examName)
	examConf := lib.GetExamConf(examName)

	// JSON 解析
	var condList map[string]string
	var sortList map[string]int
	if whereJSONStr != "" { // Note: json 不允许出现 Number 类型的 Value (eg.{"Class":1} 必须为 {"Class":"1"})
		if err := json.Unmarshal([]byte(whereJSONStr), &condList); err != nil {
			return RespError(c, "where 参数 JSON 解析失败")
		}
	}
	if sortJSONStr != "" {
		if err := json.Unmarshal([]byte(sortJSONStr), &sortList); err != nil {
			return RespError(c, "sort 参数 JSON 解析失败")
		}
	}

	// 查询条件
	if condList == nil || len(condList) == 0 {
		// 全部数据

	} else {
		if len(condList) == 1 && condList["NAME"] != "" {
			// 模糊查询
			query = lib.FilterScoresByRegStr(query, condList["NAME"])
		} else {
			// 精确查询
			query = lib.FilterScores(query, condList, false)
		}
	}

	// 数据内容描述
	dataDesc := ""
	if condList == nil {
		dataDesc = "全部考生成绩"
	} else if len(condList) == 1 && condList["NAME"] != "" {
		dataDesc = fmt.Sprintf(`数据满足 “%s” 的考生成绩`, condList["NAME"])
	} else if condList["CLASS"] == "" && condList["SCHOOL"] != "" {
		dataDesc = fmt.Sprintf(`%s · 全校成绩`, condList["SCHOOL"])
	} else if condList["CLASS"] != "" && condList["SCHOOL"] != "" {
		dataDesc = fmt.Sprintf(`%s %s · 班级成绩`, condList["SCHOOL"], condList["CLASS"])
	}

	// 排序规则
	if sortList == nil || len(sortList) == 0 {
		sortList = map[string]int{"TOTAL": -1}
	}
	for key, t := range sortList {
		if t == 1 {
			query = query.Order(key + ` asc`) // TODO: sql注入风险 待测试
		} else if t == -1 {
			query = query.Order(key + ` desc`)
		}
		break
	}

	// 分页操作
	var page, pageSize, offset, lastPage int
	var total int64
	page, _ = strconv.Atoi(pageStr)
	pageSize, _ = strconv.Atoi(pageSizeStr)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	offset = (page - 1) * pageSize
	query.Count(&total)
	lastPage = int(math.Ceil(float64(total) / float64(pageSize)))
	query = query.Offset(offset).Limit(pageSize) // 先读取完整 scoreList，再分页

	// 响应数据
	scList := []model.Score{}
	query.Find(&scList)

	fieldList := getScoresFieldList(*examConf, scList)
	//scListPaginated := scoreListPaginate(scList, offset, pageSize) // 数据分页
	scoreListAvgList := Map{} //scoreListAvgList(scList, fieldList)        // 平均分

	pageResult := Map{
		"examName":  examName,
		"dataDesc":  dataDesc,
		"page":      page,
		"pageSize":  pageSize,
		"total":     total,
		"lastPage":  lastPage,
		"fieldList": fieldList,
		"list":      scList,
		"examConf":  examConf,
		"avgList":   scoreListAvgList,
		"sortList":  sortList,
		"condList":  condList,
	}

	if initConf != nil {
		pageResult["initConf"] = initConf
	}

	return RespData(c, pageResult)
}

func scoreListPaginate(x []model.Score, skip int, size int) []model.Score {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}

func scoreListAvgList(scList []model.Score, fieldList []string) map[string]float64 {
	avgList := map[string]float64{}
	avgFields := []string{}
	avgFields = append(avgFields, model.SFieldSubj...)
	avgFields = append(avgFields, model.SFieldExtSum...)

	for _, f := range funk.IntersectString(avgFields, fieldList) {
		scores := funk.Map(scList, func(sc model.Score) float64 {
			num, err := reflections.GetField(sc, f)
			if err != nil {
				return 0
			}

			switch num := num.(type) {
			case int:
				return float64(num)
			case float64:
				return num
			default:
				return 0
			}
		}).([]float64)

		d := len(scores)
		if d == 0 {
			d = 1
		}
		avgList[f] = funk.SumFloat64(scores) / float64(d)
	}

	return avgList
}

// 获取成绩数据的可用字段
func getScoresFieldList(examConf model.ExamConf, scoreList []model.Score) (fieldList []string) {
	fieldList = []string{}
	allField, err := reflections.Fields(&model.Score{})
	if err != nil {
		return
	}

	// 将 examConf 预设的学科字段名加入
	if examConf.Subj != "" {
		var subjects []string
		if err := utils.JSONDecode(examConf.Subj, &subjects); err == nil {
			fieldList = append(fieldList, subjects...)
		}
	}

	for _, sc := range scoreList {
		for _, fn := range allField {
			val, err := reflections.GetField(sc, fn)
			if err != nil {
				continue
			}
			nullVal, _ := reflections.GetField(model.Score{}, fn)
			if val != nullVal && !funk.ContainsString(fieldList, fn) {
				fieldList = append(fieldList, fn)
			}
		}
	}

	return fieldList
}
