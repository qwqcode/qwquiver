package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// Map is a map
	Map = map[string]interface{}
)

// JSONResult JSON 响应数据结构
type JSONResult struct {
	Success bool        `json:"success"` // 是否成功
	Msg     string      `json:"msg"`     // 消息
	Data    interface{} `json:"data"`    // 数据
	Extra   interface{} `json:"extra"`   // 数据
}

// RespJSON is normal json result
func RespJSON(c echo.Context, msg string, data interface{}, success bool) error {
	return c.JSON(http.StatusOK, &JSONResult{
		Success: success,
		Msg:     msg,
		Data:    data,
	})
}

// RespData is just response data
func RespData(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, &JSONResult{
		Success: true,
		Data:    data,
	})
}

// RespSuccess is just response success
func RespSuccess(c echo.Context) error {
	return c.JSON(http.StatusOK, &JSONResult{
		Success: true,
	})
}

// RespError is just response error
func RespError(c echo.Context, msg string, details ...string) error {
	return c.JSON(http.StatusInternalServerError, &JSONResult{
		Success: false,
		Msg:     msg,
		Extra: Map{
			"errDetails": details,
		},
	})
}
