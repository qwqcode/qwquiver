package http

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
)

// GetAll Api: /api/school/all
func schoolAllHandler(c echo.Context) error {
	examName := c.QueryParam("exam")

	if !lib.IsExamExist(examName) {
		return RespError(c, "Exam 不存在")
	}
	exam := lib.GetExam(examName)
	// examConf := lib.GetExamConf(examName)

	schoolList := map[string][]string{}

	exam.Select().Each(new(model.Score), func(record interface{}) error {
		sc := record.(*model.Score)
		school := strings.TrimSpace(sc.SCHOOL)
		if school != "" {
			if schoolList[school] == nil {
				schoolList[school] = []string{}
			}
			class := strings.TrimSpace(sc.CLASS)
			if !funk.ContainsString(schoolList[school], class) {
				schoolList[school] = append(schoolList[school], class)
			}
		}
		return nil
	})

	return RespData(c, Map{
		"school": schoolList,
	})
}
