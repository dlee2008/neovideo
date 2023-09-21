package jiexi

import (
	"fmt"

	"d1y.io/neovideo/common/impl"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/kataras/iris/v12"
)

type JiexiController struct {
}

func (jx *JiexiController) getList(ctx iris.Context) {
	jiexi, db := gplus.SelectList[repos.JiexiRepoWithModel](nil)
	if db.Error != nil {
		web.NewJSONResultWithError(db.Error).Build(ctx)
		return
	}
	web.NewJSONResultWithSuccess(jiexi).Build(ctx)
}

func (jx *JiexiController) create(ctx iris.Context) {
	var jiexiForm repos.JiexiRepo
	ctx.ReadBody(&jiexiForm)
	var insertData = repos.JiexiRepoWithModel{
		JiexiRepo: jiexiForm,
	}
	err := gplus.Insert[repos.JiexiRepoWithModel](&insertData).Error
	if err != nil {
		web.NewJSONResultWithError(err).Build(ctx)
		return
	}
	web.NewJSONResultWithSuccess(insertData)
}

func (jx *JiexiController) delete(ctx iris.Context) {
	id := ctx.Params().Get("id")
	err := gplus.DeleteById[repos.JiexiRepoWithModel](id).Error
	if err != nil {
		web.NewJSONResultWithError(err).Build(ctx)
		return
	}
	web.NewJSONResultWithSuccess(id).SetMessage("删除成功").Build(ctx)
}

func (jx *JiexiController) batchImport(ctx iris.Context) {
	importData := ctx.FormValueDefault("data", "")
	if len(importData) == 0 {
		web.NewJSONResultWithMessage("导入数据为空").SetSuccessWithBool(false).Build(ctx)
		return
	}
	list := impl.ParseJiexi(importData)
	if len(list) == 0 {
		web.NewJSONResultWithMessage("导入数据为空").SetSuccessWithBool(false).Build(ctx)
		return
	}
	var importJiexi = make([]*repos.JiexiRepoWithModel, 0)
	for _, item := range list {
		importJiexi = append(importJiexi, &repos.JiexiRepoWithModel{
			JiexiRepo: repos.JiexiRepo{
				Name: item.Name,
				URL:  item.URL,
			},
		})
	}
	err := gplus.InsertBatch(importJiexi).Error
	if err != nil {
		web.NewJSONResultWithError(err).Build(ctx)
		return
	}
	web.NewJSONResultWithMessage(fmt.Sprintf("新增成功(%d条)", len(importJiexi))).SetSuccessWithBool(true).Build(ctx)
}

func Register(u iris.Party) {
	var jx JiexiController
	u.Get("/", jx.getList)
	u.Post("/", jx.create)
	u.Delete("/{id:int}", jx.delete)
	u.Post("/batch_import", jx.batchImport)
}