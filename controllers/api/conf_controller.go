package api

import (
	"github.com/kataras/iris/v12"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
)

// ConfController 配置数据控制器
type ConfController struct {
	Ctx iris.Context
}

// Get Api: /api/conf [已废弃]
func (c *ConfController) Get() *utils.JSONResult {
	examList := lib.GetAllExamConf()
	examGrpList := lib.GetAllExamGrps()
	fieldTransDict := model.ScoreFieldTransMap

	return utils.JSONData(iris.Map{
		"examList":       examList,
		"examGrpList":    examGrpList,
		"fieldTransDict": fieldTransDict,
	})
}
