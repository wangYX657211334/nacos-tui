package config

import (
	"fmt"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	nacos "github.com/wangYX657211334/nacos-tui/pkg/nacos"
	"github.com/wangYX657211334/nacos-tui/pkg/util"
)

type NacosConfigEdit struct {
	repo repository.Repository
}

func (m *NacosConfigEdit) EditConfigContent(dataId string, group string, editCallBack func() error) error {
	res, err := m.repo.UpdateConfig(dataId, group, func(ncr *nacos.ConfigResponse) (res bool, err error) {
		res, ncr.Content, err = util.EditContentBySystemEditor(fmt.Sprintf("%s@%s", group, dataId), ncr.Content)
		return
	})
	if err != nil {
		return err
	}
	if res {
		err := editCallBack()
		if err != nil {
			return err
		}
		event.Publish(event.ApplicationMessageEvent, "配置已更新")
	} else {
		event.Publish(event.ApplicationMessageEvent, "未修改，无需更新")
	}
	return nil
}
