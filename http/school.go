package http

import (
	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
)

// GetAll Api: /api/school/all
func schoolAllHandler(c echo.Context) error {
	examName := c.QueryParam("exam")

	if !lib.HasExam(examName) {
		return RespError(c, "Exam 不存在")
	}
	// examConf := lib.GetExamConf(examName)

	result := map[string][]string{}

	schoolList := []string{}
	lib.NewExamQuery(examName).Select("school").Group("school").Find(&schoolList)

	for _, school := range schoolList {
		if school != "" {
			if result[school] == nil {
				result[school] = []string{}
			}

			classList := []string{}
			lib.NewExamQuery(examName).Select("class").Where("school = ?", school).Group("class").Find(&classList)
			result[school] = append(result[school], classList...)
		}
	}

	return RespData(c, Map{
		"school": result,
	})
}
