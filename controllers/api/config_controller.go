package api

import (
	"github.com/kataras/iris/v12"
	"github.com/qwqcode/qwquiver/lib/utils"
)

// ConfigController 配置数据控制器
type ConfigController struct {
	Ctx iris.Context
}

// Get Api: /api/config
func (c *ConfigController) Get() *utils.JSONResult {
	return utils.JSONData(map[string]string{"hello": "jsonp"})
}
