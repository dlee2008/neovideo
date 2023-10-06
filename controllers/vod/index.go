package vod

import (
	"sync"
	"time"

	"d1y.io/neovideo/controllers/handler"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"github.com/kataras/iris/v12"
	"github.com/patrickmn/go-cache"
)

type IVodController struct {
	sm sync.Mutex
	cc *cache.Cache
}

func newVod() *IVodController {
	vod := IVodController{}
	vod.cc = cache.New(42*time.Second, 60*time.Second)
	return &vod
}

func (ic *IVodController) getVideos(ctx iris.Context) {
	data, err := handler.BuildPagination[repos.VideoRepo](ctx)
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	web.NewData(data).Build(ctx)
}

func Register(u iris.Party) {
	vod := newVod()

	// Deprecated: remove this
	u.Get("/home", vod.renderHome).Name = "代理访问首页"

	u.Get("/videos", vod.getVideos)
}
