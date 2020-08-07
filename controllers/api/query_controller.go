package api

import (
	"github.com/kataras/iris/v12"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
)

type QueryController struct {
	Ctx iris.Context
}

func (c *QueryController) Get() *utils.JSONResult {
	classes := []string{}
	lib.FilterScoresRegStr("233", "林").Each(new(model.Score), func(record interface{}) error {
		u := record.(*model.Score)
		classes = append(classes, u.Name)
		return nil
	})

	return utils.JSONData(map[string]interface{}{"data": classes})
}
