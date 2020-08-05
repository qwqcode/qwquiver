package api

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/qwqcode/qwquiver/lib"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"go.etcd.io/bbolt"
)

type QueryController struct {
	Ctx iris.Context
}

func (c *QueryController) Get() *utils.JSONResult {
	lib.DB.From("233s").Init(&model.Score{})
	s := &[]string{}
	fmt.Printf("%p\n", s)

	_ = lib.DB.Bolt.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			*s = append(*s, string(name))
			return nil
		})
	})

	fmt.Printf("%p\n", s)

	return utils.JSONData(map[string]interface{}{"data": *s, "p1": s})
}
