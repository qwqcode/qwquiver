package http

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
	"gopkg.in/oleiade/reflections.v1"
	"gorm.io/gorm"
)

// query API 的公共参数
type queryAPICommonParms struct {
	examName     string
	whereJSONStr string
	pageStr      string
	pageSizeStr  string
	sortJSONStr  string
	initConf     Map

	query       *gorm.DB
	examConf    *model.ExamConf
	condList    map[string]string
	sortList    map[string]int
	dataDesc    string
	subjectList []string
}

func getQueryAPICommonParms(c echo.Context) *queryAPICommonParms {
	p := &queryAPICommonParms{
		examName:     c.QueryParam("exam"),
		whereJSONStr: c.QueryParam("where"),
		pageStr:      c.QueryParam("page"),
		pageSizeStr:  c.QueryParam("pageSize"),
		sortJSONStr:  c.QueryParam("sort"),
	}

	isInitReq := c.QueryParam("init") != ""
	if isInitReq {
		// 若为初始化请求
		examMap := lib.GetAllExamsSorted()
		examGrpList := lib.GetAllExamGrps()
		fieldTransDict := model.ScoreFieldTransMap

		if len(examMap) == 0 {
			RespError(c, "未找到任何考试数据，请导入数据")
			return nil
		}

		if p.examName == "" {
			p.examName = lib.GetAllExamsSorted()[0].Name // 设置默认 exam
		}

		p.initConf = Map{
			"examMap":        examMap,
			"examGrpList":    examGrpList,
			"fieldTransDict": fieldTransDict,
		}
	}

	if !lib.HasExam(p.examName) {
		RespError(c, "Exam 不存在")
		return nil
	}

	p.query = lib.NewExamQuery(p.examName)
	p.examConf = lib.GetExamConf(p.examName)

	// JSON 解析
	if p.whereJSONStr != "" { // Note: json 不允许出现 Number 类型的 Value (eg.{"Class":1} 必须为 {"Class":"1"})
		if err := json.Unmarshal([]byte(p.whereJSONStr), &p.condList); err != nil {
			RespError(c, "where 参数 JSON 解析失败", err.Error())
			return nil
		}
	}
	if p.sortJSONStr != "" {
		if err := json.Unmarshal([]byte(p.sortJSONStr), &p.sortList); err != nil {
			RespError(c, "sort 参数 JSON 解析失败", err.Error())
			return nil
		}
	}

	// 数据内容描述
	if p.condList == nil {
		p.dataDesc = "全部考生成绩"
	} else if len(p.condList) == 1 && p.condList["NAME"] != "" {
		p.dataDesc = fmt.Sprintf(`数据满足 “%s” 的考生成绩`, p.condList["NAME"])
	} else if p.condList["CLASS"] == "" && p.condList["SCHOOL"] != "" {
		p.dataDesc = fmt.Sprintf(`%s · 全校成绩`, p.condList["SCHOOL"])
	} else if p.condList["CLASS"] != "" && p.condList["SCHOOL"] != "" {
		p.dataDesc = fmt.Sprintf(`%s %s · 班级成绩`, p.condList["SCHOOL"], p.condList["CLASS"])
	}

	// 查询条件
	if p.condList == nil || len(p.condList) == 0 {
		// 全部数据

	} else {
		if len(p.condList) == 1 && p.condList["NAME"] != "" {
			// 模糊查询
			p.query = lib.FilterScoresByRegStr(p.query, p.condList["NAME"])
		} else {
			// 精确查询
			p.query = lib.FilterScores(p.query, p.condList, false)
		}
	}

	// 排序规则
	if p.sortList == nil || len(p.sortList) == 0 {
		p.sortList = map[string]int{"TOTAL": -1}
	}
	for key, t := range p.sortList {
		if t == 1 {
			p.query = p.query.Order(key + ` asc`) // TODO: sql注入风险 待测试
		} else if t == -1 {
			p.query = p.query.Order(key + ` desc`)
		}
		break
	}

	// 考试科目
	var err error
	p.subjectList, err = getExamSubjectList(*p.examConf)
	if err != nil {
		RespError(c, "考试科目数据获取失败", err.Error())
		return nil
	}

	return p
}

// 查询 get api/query
func queryHandler(c echo.Context) error {
	p := getQueryAPICommonParms(c)
	if p == nil {
		return nil
	}

	// 分页操作
	var total int64 // 数据总数
	p.query.Count(&total)

	var page, pageSize int // 页码
	page, _ = strconv.Atoi(p.pageStr)
	if page <= 0 {
		page = 1
	}
	pageSize, _ = strconv.Atoi(p.pageSizeStr)
	if pageSize <= 0 {
		pageSize = 50
	}

	var lastPage int // 最后一页
	var offset int   // 数据偏移量
	lastPage = int(math.Ceil(float64(total) / float64(pageSize)))
	offset = (page - 1) * pageSize

	// 响应数据
	scList := []model.Score{}
	p.query.Offset(offset).Limit(pageSize).Find(&scList)

	pageResult := Map{
		"examName":    p.examName,
		"dataDesc":    p.dataDesc,
		"page":        page,
		"pageSize":    pageSize,
		"total":       total,
		"lastPage":    lastPage,
		"subjectList": p.subjectList,
		"list":        scList,
		"examConf":    p.examConf,
		"sortList":    p.sortList,
		"condList":    p.condList,
	}

	if p.initConf != nil {
		pageResult["initConf"] = p.initConf
	}

	return RespData(c, pageResult)
}

// 求均值 get api/query/avg
func queryAvgHandler(c echo.Context) error {
	p := getQueryAPICommonParms(c)
	if p == nil {
		return nil
	}

	avgKeys := []string{"TOTAL", "LZ", "WZ", "ZK"}
	avgKeys = append(avgKeys, p.subjectList...)

	avgSQL := ""
	for _, k := range avgKeys {
		avgSQL += "avg(" + k + ") as " + k + ","
	}
	avgSQL = strings.TrimRight(avgSQL, ",")

	var result []map[string]interface{}
	p.query.Select(avgSQL).Find(&result)

	avgList := result[0]
	return RespData(c, Map{
		"avgList": avgList,
	})
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

// 获取 exam 的所有科目
func getExamSubjectList(examConf model.ExamConf) (subjects []string, err error) {
	if examConf.Subj != "" {
		if err := utils.JSONDecode(examConf.Subj, &subjects); err != nil {
			return []string{}, err
		}
	}
	return
}

// 获取成绩数据的可用字段
func getScoresFieldList(examConf model.ExamConf, scoreSample model.Score) (fieldList []string) {
	fieldList = []string{}

	// 将 examConf 预设的学科字段名加入
	if examConf.Subj != "" {
		var subjects []string
		if err := utils.JSONDecode(examConf.Subj, &subjects); err == nil {
			fieldList = append(fieldList, subjects...)
		}
	}

	allField, err := reflections.Fields(&model.Score{})
	if err != nil {
		return
	}
	for _, fn := range allField {
		val, err := reflections.GetField(scoreSample, fn)
		if err != nil {
			continue
		}
		nullVal, _ := reflections.GetField(model.Score{}, fn)
		if val != nullVal && !funk.ContainsString(fieldList, fn) {
			fieldList = append(fieldList, fn)
		}
	}

	return fieldList
}
