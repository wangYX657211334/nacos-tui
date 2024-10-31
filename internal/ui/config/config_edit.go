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

func (m *NacosConfigEdit) EditConfigContent(command string, dataId string, group string) (tea.Cmd, error) {
	configRes, err := m.repo.GetConfig(dataId, group)
	if err != nil {
		return nil, err
	}
	return util.EditContent(command, fmt.Sprintf("%s@%s", group, dataId), configRes.Content, func(ok bool, newContentYaml string, err error) {
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
				event.Publish(event.ApplicationMessageEvent, "配置已更新")
			} else {
				event.Publish(event.ApplicationMessageEvent, "配置更新失败")
			}
		} else {
			event.Publish(event.ApplicationMessageEvent, "未修改，无需更新")
		}
	}), nil
}
