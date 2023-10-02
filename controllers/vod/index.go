package vod

import (
	"errors"
	"sort"
	"sync"
	"time"

	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/kataras/iris/v12"
	"github.com/patrickmn/go-cache"
)

const (
	homeRepoQueryKey = "$home"
	homeRenderKey    = "$hrender"
)

type homeItem struct {
	ID    uint                   `json:"id"`
	Api   string                 `json:"api"`
	Name  string                 `json:"name"`
	Data  maccms.IMacCMSHomeData `json:"data,omitempty"`
	Error string                 `json:"error,omitempty"`
}

type IVodController struct {
	sm sync.Mutex
	cc *cache.Cache
}

func newVod() *IVodController {
	vod := IVodController{}
	vod.cc = cache.New(42*time.Second, 60*time.Second)
	return &vod
}

func (vc *IVodController) queryRawCMS() ([]repos.MacCMSRepo, error) {
	vc.sm.Lock()
	defer vc.sm.Unlock()
	var result []repos.MacCMSRepo
	if val, ok := vc.cc.Get(homeRepoQueryKey); ok {
		if v, o := val.([]repos.MacCMSRepo); o {
			result = v
		}
	} else {
		cms, gb := gplus.SelectList[repos.MacCMSRepo](nil)
		if gb.Error != nil {
			return nil, gb.Error
		}
		for _, item := range cms {
			result = append(result, *item)
		}
		vc.cc.SetDefault(homeRepoQueryKey, result)
	}
	return result, nil
}

func (vc *IVodController) queryAndCMSFetchHome() ([]homeItem, error) {
	c, e := vc.queryRawCMS()
	if e != nil {
		return nil, e
	}
	if len(c) <= 0 {
		return nil, errors.New("maccms is empty")
	}
	var wg sync.WaitGroup
	var data []homeItem
	wg.Add(len(c))
	for _, item := range c {
		go func(item repos.MacCMSRepo) {
			defer wg.Done()
			var im = homeItem{
				Name: item.Name,
				ID:   item.ID,
				Api:  item.Api,
			}
			val, err := maccms.New(item.RespType, item.Api).GetHome()
			if err != nil {
				im.Error = err.Error()
			} else {
				im.Data = val
			}
			data = append(data, im)
		}(item)
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].ID < data[j].ID
	})
	wg.Wait()
	return data, nil
}

func (vc *IVodController) renderHome(ctx iris.Context) {
	var result []homeItem
	if val, ok := vc.cc.Get(homeRenderKey); ok {
		if v, o := val.([]homeItem); o {
			result = v
		}
	} else {
		ims, err := vc.queryAndCMSFetchHome()
		if err != nil {
			web.NewError(err).Build(ctx)
			return
		}
		vc.cc.SetDefault(homeRenderKey, ims)
		result = ims
	}
	web.NewData(result).Build(ctx)
}

func Register(u iris.Party) {
	vod := newVod()
	u.Get("/home", vod.renderHome)
}
