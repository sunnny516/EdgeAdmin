package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	webConfig, err := webutils.FindWebConfigWithLocationId(this.Parent(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["cacheConfig"] = webConfig.Cache

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId     int64
	CacheJSON []byte

	Must *actions.Must
}) {
	// TODO 校验配置
	cacheConfig := &serverconfigs.HTTPCacheConfig{}
	err := json.Unmarshal(params.CacheJSON, cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 去除不必要的部分
	for _, cacheRef := range cacheConfig.CacheRefs {
		cacheRef.CachePolicy = nil
	}

	cacheJSON, err := json.Marshal(cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebCache(this.AdminContext(), &pb.UpdateHTTPWebCacheRequest{
		WebId:     params.WebId,
		CacheJSON: cacheJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}