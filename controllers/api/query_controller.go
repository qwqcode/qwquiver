package api

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/asdine/storm"
	"github.com/kataras/iris/v12"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
)

type QueryController struct {
	Ctx iris.Context
}

func (c *QueryController) Get() *utils.JSONResult {
	examName := c.Ctx.URLParamDefault("exam", "")
	whereJsonStr := c.Ctx.URLParamDefault("where", "")
	pageStr := c.Ctx.URLParamDefault("page", "")
	pageSizeStr := c.Ctx.URLParamDefault("pageSize", "")
	sortJsonStr := c.Ctx.URLParamDefault("sort", "")

	if !lib.IsExamExist(examName) {
		return utils.JSONError(utils.RespCodeErr, "Exam 不存在")
	}
	exam := lib.GetExam(examName)

	var condList map[string]interface{}
	var sortList map[string]int

	// JSON 解析
	if whereJsonStr != "" {
		if err := json.Unmarshal([]byte(whereJsonStr), &condList); err != nil {
			return utils.JSONError(utils.RespCodeErr, "where 参数 JSON 解析失败")
		}
	}
	if sortJsonStr != "" {
		if err := json.Unmarshal([]byte(sortJsonStr), &sortList); err != nil {
			return utils.JSONError(utils.RespCodeErr, "sort 参数 JSON 解析失败")
		}
	}

	if sortList == nil || len(sortList) == 0 {
		sortList = map[string]int{"Total": -1}
	}

	var page int
	var pageSize int
	page, _ = strconv.Atoi(pageStr)
	pageSize, _ = strconv.Atoi(pageSizeStr)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var query storm.Query
	if condList == nil || len(condList) == 0 {
		// 全部数据
		query = exam.Select()
	} else {
		if len(condList) == 1 && condList["Name"] != "" {
			// 模糊查询
			query = lib.FilterScoresByRegStr(exam, condList["Name"].(string))
		} else {
			query = lib.FilterScores(exam, condList, false)
		}
	}

	total, _ := query.Count(&model.Score{})
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))

	// 排序
	for key, t := range sortList {
		query = query.OrderBy(key)
		if t == -1 {
			query = query.Reverse()
		}
		break
	}

	query = query.Skip(offset).Limit(pageSize)

	scList := []model.Score{}
	query.Each(new(model.Score), func(record interface{}) error {
		sc := record.(*model.Score)
		scList = append(scList, *sc)
		return nil
	})

	pageResult := iris.Map{
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
		"lastPage": lastPage,
		"list":     scList,
	}

	return utils.JSONData(pageResult)
}
