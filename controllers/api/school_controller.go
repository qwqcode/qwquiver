package api

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/thoas/go-funk"
)

// SchoolController 配置数据控制器
type SchoolController struct {
	Ctx iris.Context
}

// GetAll Api: /api/school/all
func (c *SchoolController) GetAll() *utils.JSONResult {
	examName := c.Ctx.URLParamDefault("exam", "")

	if !lib.IsExamExist(examName) {
		return utils.JSONError(utils.RespCodeErr, "Exam 不存在")
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

	return utils.JSONData(iris.Map{
		"school": schoolList,
	})
}
