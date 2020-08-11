package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/model"
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
		var queryPersonSc []model.Score
		if err := lib.FilterScores(exam.Data, condList, false).Find(&queryPersonSc); err != nil {
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
		if exam.Conf.Subj == nil || len(exam.Conf.Subj) == 0 {
			continue
		}
		subjScores := map[string]float64{}
		for _, f := range exam.Conf.Subj {
			scoreI, err := reflections.GetField(sc, f)
			if err != nil {
				continue
			}
			score := float64(scoreI.(float64))
			if exam.Conf.SubjFullScore != nil && exam.Conf.SubjFullScore[f] > 0 {
				score = (score / exam.Conf.SubjFullScore[f]) * 100            // 转为百分制
				score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", score), 64) // 保留两位小数
			}

			subjScores[f] = score
		}

		chartData = append(chartData, subjScores)
	}

	return RespData(c, Map{
		"chartData": chartData,
		"fieldList": fields,
		"uncertain": uncertain,
	})
}
