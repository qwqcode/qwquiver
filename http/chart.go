package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/oleiade/reflections.v1"
)

// Get Api: /api/chart
func chartHandler(c echo.Context) error {
	examGrp := c.QueryParam("examGrp")

	whereJSONStr := c.QueryParam("where")
	var condList map[string]string
	if whereJSONStr != "" {
		if err := json.Unmarshal([]byte(whereJSONStr), &condList); err != nil {
			return RespError(c, "where 参数 JSON 解析失败")
		}
	}

	uncertain := false // 是否为不确定的数据
	examList := lib.GetExamsByGrp(examGrp)
	fields := []string{}
	fields = append(fields, model.SFieldSubj...)

	chartData := []interface{}{}
	for _, exam := range examList {
		queryPersonSc := []model.Score{}
		if rs := lib.FilterScores(lib.NewExamQuery(exam.Name), condList, false).Find(&queryPersonSc); rs.Error != nil {
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

		// 统计此人此次考试的各科分数

		var subjects []string
		if exam.Conf.Subj != "" {
			if err := utils.JSONDecode(exam.Conf.Subj, &subjects); err != nil {
				continue
			}
			if len(subjects) == 0 {
				continue
			}
		} else {
			continue
		}

		var subjFullScore map[string]float64
		if err := utils.JSONDecode(exam.Conf.SubjFullScore, &subjFullScore); err != nil {
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

			subjScores[f] = score
		}

		examKey := exam.Conf.Label
		if examKey == "" {
			examKey = exam.Conf.Name
		}
		subjScores["exam"] = examKey
		chartData = append(chartData, subjScores)
	}

	return RespData(c, Map{
		"chartData": chartData,
		"fieldList": fields,
		"uncertain": uncertain,
	})
}
