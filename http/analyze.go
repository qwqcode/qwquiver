package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib/exd"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gopkg.in/oleiade/reflections.v1"
)

// Get Api: /api/analyze
func analyzeHandler(c echo.Context) error {
	examGrp := c.QueryParam("examGrp")

	whereJSONStr := c.QueryParam("where")
	var condList map[string]string
	if whereJSONStr != "" {
		if err := json.Unmarshal([]byte(whereJSONStr), &condList); err != nil {
			return RespError(c, "where 参数 JSON 解析失败")
		}
	}

	uncertain := false // 是否为不确定的数据
	examList := exd.GetExamsByGrp(examGrp)
	fields := []string{}
	for _, f := range model.SFieldSubj {
		fields = append(fields, model.ScoreFieldTransMap[f])
	}

	classList := []string{} // 从数据中收集该姓名存在的班级

	// 获取每次考试的数据
	exams := []interface{}{}
	{
		for _, exam := range examList {
			queryPersonSc := []model.Score{}
			if rs := exd.FilterScores(exam.NewQuery(), condList, false).Find(&queryPersonSc); rs.Error != nil {
				logrus.Error("api.chart ", rs.Error)
				continue
			}
			if len(queryPersonSc) == 0 {
				continue
			}
			if len(queryPersonSc) > 1 {
				uncertain = true
			}

			sc := queryPersonSc[0]

			// 记录该姓名存在的班级
			if !funk.ContainsString(classList, sc.CLASS) {
				classList = append(classList, sc.CLASS)
			}

			// 统计此人此次考试的各科分数
			var subjects []string
			if exam.Subj != "" {
				if err := utils.JSONDecode(exam.Subj, &subjects); err != nil {
					continue
				}
				if len(subjects) == 0 {
					continue
				}
			} else {
				continue
			}

			var subjFullScore map[string]float64
			if err := utils.JSONDecode(exam.SubjFullScore, &subjFullScore); err != nil {
				continue
			}
			subjScores := map[string]interface{}{}
			for _, f := range subjects {
				scoreI, err := reflections.GetField(sc, f)
				if err != nil {
					continue
				}
				score := float64(scoreI.(float64))
				if subjFullScore != nil && subjFullScore[f] > 0 {
					score = (score / subjFullScore[f]) * 100                      // 转为百分制
					score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", score), 64) // 保留两位小数
				}

				subjScores[model.ScoreFieldTransMap[f]] = score
			}

			examKey := exam.Label
			if examKey == "" {
				examKey = exam.Name
			}
			subjScores["exam"] = examKey
			subjScores["date"] = exam.Date
			exams = append(exams, subjScores)
		}
	}

	result := Map{
		"examGrp":   examGrp,
		"name":      condList["NAME"],
		"school":    condList["SCHOOL"],
		"classList": classList,
		"exams":     exams,
		"examCount": len(exams),
		"fieldList": fields,
		"uncertain": uncertain,
	}

	return RespData(c, result)
}
