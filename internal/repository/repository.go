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

func NewRepository(db *sql.DB) (Repository, error) {
	r := &repository{}
	appConfig, err := config.LoadApplicationConfig()
	if err != nil {
		return nil, err
	}
	r.configApi = appConfig
	r.nacosApi = nacos.NewApi(func() (url string, username string, password string, namespace string) {
		ctx := appConfig.GetNacosContext()
		return ctx.Url, ctx.User, ctx.Password, ctx.UseNamespace
	}, db)
	return r, nil
}
