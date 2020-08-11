package http

import (
	"github.com/labstack/echo/v4"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/model"
)

// Get Api: /api/conf [已废弃]
func confHandler(c echo.Context) error {
	examList := lib.GetAllExamConf()
	examGrpList := lib.GetAllExamGrps()
	fieldTransDict := model.ScoreFieldTransMap

	return RespData(c, Map{
		"examList":       examList,
		"examGrpList":    examGrpList,
		"fieldTransDict": fieldTransDict,
	})
}
