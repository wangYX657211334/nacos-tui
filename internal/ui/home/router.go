package home

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strings"

	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/internal/ui/config"
	"github.com/wangYX657211334/nacos-tui/internal/ui/context"
	"github.com/wangYX657211334/nacos-tui/internal/ui/namespace"
	"github.com/wangYX657211334/nacos-tui/internal/ui/service"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type Router struct {
	Name          string
	Path          string
	RootComponent bool
	Component     func(repository.Repository, ...any) NacosModel
}

var DefaultRoutePath = "/context"

func DefaultRoute() (r Router) {
	for _, router := range Routers {
		if strings.EqualFold(router.Path, DefaultRoutePath) {
			r = router
			return
		}
	}
	return
}

var Routers = []Router{
	{
		Name:          "context",
		Path:          "/context",
		RootComponent: true,
		Component: func(repo repository.Repository, param ...any) NacosModel {
			return context.NewNacosContextModel(repo)
		},
	},
	{
		Name:          "namespace",
		Path:          "/namespace",
		RootComponent: true,
		Component: func(repo repository.Repository, param ...any) NacosModel {
			return namespace.NewNacosNamespaceModel(repo)
		},
	},
	{
		Name:          "config",
		Path:          "/config",
		RootComponent: true,
		Component: func(repo repository.Repository, param ...any) NacosModel {
			return config.NewNacosConfigListModel(repo)
		},
	},
	{
		Name: "listener",
		Path: "/config/listener",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			dataId := param[0].(string)
			group := param[1].(string)
			return config.NewNacosConfigListenerModel(repo, dataId, group)
		},
	},
	{
		Name: "history",
		Path: "/config/history",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			dataId := param[0].(string)
			group := param[1].(string)
			return config.NewNacosConfigHistoryModel(repo, dataId, group)
		},
	},
	{
		Name: "confirm(clone)",
		Path: "/config/clone",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			configs := param[0].([]nacos.ConfigsItem)
			return config.NewNacosConfigClone(repo, configs)
		},
	},
	{
		Name: "result(clone)",
		Path: "/config/clone/result",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			configs := param[0].(*nacos.ConfigCloneResponse)
			return config.NewNacosConfigCloneResultModel(repo, configs)
		},
	},
	{
		Name: "confirm(delete)",
		Path: "/config/delete",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			configs := param[0].([]nacos.ConfigsItem)
			return config.NewNacosConfigDeleteModel(repo, configs)
		},
	},
	{
		Name:          "service",
		Path:          "/service",
		RootComponent: true,
		Component: func(repo repository.Repository, param ...any) NacosModel {
			return service.NewNacosServiceListModel(repo)
		},
	},
	{
		Name: "instance",
		Path: "/service/instance",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			dataId := param[0].(string)
			group := param[1].(string)
			return service.NewNacosServiceInstanceListModel(repo, dataId, group)
		},
	},
	{
		Name: "subscriber",
		Path: "/service/subscriber",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			dataId := param[0].(string)
			group := param[1].(string)
			return service.NewNacosServiceSubscriberListModel(repo, dataId, group)
		},
	},
	{
		Name: "command",
		Path: "/**/command",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			content := param[0].(base.Model)
			suggestions := param[1].([]base.Suggestion)
			execute := param[2].(func(string) bool)
			return NewCommandModel(content, suggestions, execute)
		},
	},
	{
		Name: "filter",
		Path: "/**/filter",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			defaultDataId := param[0].(string)
			defaultGroup := param[1].(string)
			content := param[2].(base.Model)
			filter := param[3].(func(dataId, group string))
			return base.NewListFilterModel(repo, defaultDataId, defaultGroup, content, filter)
		},
	},
	{
		Name: "view",
		Path: "/**/view",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			content := param[0].(string)
			return base.NewDetailModel(repo, content)
		},
	},
	{
		Name: "select",
		Path: "/**/select",
		Component: func(repo repository.Repository, param ...any) NacosModel {
			items := param[0].([]base.SelectItem)
			selectHandle := param[1].(func(item *base.SelectItem))
			return base.NewSelectListModel(repo, items, selectHandle)
		},
	},
}
