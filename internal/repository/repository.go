package repository

import (
	"database/sql"
	"github.com/wangYX657211334/nacos-tui/internal/config"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
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
	r.nacosApi = nacos.NewApi(func() (url string, username string, password string, namespace string, err error) {
		ctx, err := appConfig.GetNacosContext()
		if err != nil {
			return "", "", "", "", err
		}
		return ctx.Url, ctx.User, ctx.Password, ctx.UseNamespace, nil
	}, db)
	return r
}
