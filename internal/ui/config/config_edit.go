package config

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/util"
)

type NacosConfigEdit struct {
	repo repository.Repository
}

func (m *NacosConfigEdit) EditConfigContent(dataId string, group string, editCallBack func() error) (tea.Cmd, error) {
	configRes, err := m.repo.GetConfig(dataId, group)
	if err != nil {
		return nil, err
	}
	return util.EditContentBySystemEditor(fmt.Sprintf("%s@%s", group, dataId), configRes.Content, func(ok bool, newContentYaml string, err error) {
		if err != nil {
			event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
			return
		}
		if ok {
			configRes.Content = newContentYaml
			res, err := m.repo.UpdateConfig(configRes)
			if err != nil {
				event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
			}
			if res {
				err := editCallBack()
				if err != nil {
					event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
				}
				event.Publish(event.ApplicationMessageEvent, "配置已更新")
			} else {
				event.Publish(event.ApplicationMessageEvent, "未修改，无需更新")
			}
		}
	}), nil
}
