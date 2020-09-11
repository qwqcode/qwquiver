package http

import (
	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib/exd"
)

// GetAll Api: /api/school/all
func schoolAllHandler(c echo.Context) error {
	examName := c.QueryParam("exam")
	exam := exd.GetExam(examName)

	if exam == nil {
		return RespError(c, "Exam 不存在")
	}
	// examConf := lib.GetExamConf(examName)

	result := map[string][]string{}

	schoolList := []string{}
	exam.NewQuery().Select("school").Group("school").Find(&schoolList)

	for _, school := range schoolList {
		if school != "" {
			if result[school] == nil {
				result[school] = []string{}
			}

			classList := []string{}
			exam.NewQuery().Select("class").Where("school = ?", school).Group("class").Find(&classList)
			result[school] = append(result[school], classList...)
		}
	}

	return RespData(c, Map{
		"school": result,
	})
}
