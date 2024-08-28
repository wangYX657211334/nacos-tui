package repository

import (
	"database/sql"
	"github.com/wangYX657211334/nacos-tui/internal/config"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
	"log"
)

type Repository interface {
	nacos.Api
	config.Api
}

type nacosApi nacos.Api
type configApi config.Api
type repository struct {
	nacosApi
	configApi
}

func NewRepository(db *sql.DB) Repository {
	r := &repository{}
	appConfig := config.NewApi(db)
	r.configApi = appConfig
	r.nacosApi = nacos.NewApi(func() (url string, username string, password string, namespace string) {
		ctx, err := appConfig.GetNacosContext()
		if err != nil {
			log.Panic(err)
		}
		return ctx.Url, ctx.User, ctx.Password, ctx.UseNamespace
	}, db)
	return r
}
