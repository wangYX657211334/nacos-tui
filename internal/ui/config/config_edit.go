package config

import (
	"errors"
	"fmt"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	nacos "github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type NacosConfigEdit struct {
	repo repository.Repository
}

func (m *NacosConfigEdit) EditConfigContent(dataId string, group string, editCallBack func() error) error {
	res, err := m.repo.UpdateConfig(dataId, group, func(ncr *nacos.ConfigResponse) (res bool, err error) {
		fileName := filepath.Join(os.TempDir(), "nacos-tui", fmt.Sprintf("%s@%s", group, dataId))
		err = os.MkdirAll(filepath.Dir(fileName), 0755)
		if err != nil {
			return
		}
		if err := os.WriteFile(fileName, []byte(ncr.Content), 0644); err != nil {
			log.Panic(err)
		}
		defer func(name string) {
			rErr := os.Remove(name)
			if rErr != nil {
				err = errors.Join(errors.New("remove tmp file error"), rErr, err)
			}
			return
		}(fileName)
		cmd := exec.Command("vim", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return false, errors.Join(errors.New("edit file error"), err)
		}
		newFileContentBytes, err := os.ReadFile(fileName)
		if err != nil {
			return false, errors.Join(errors.New("read tmp file error"), err)
		}
		oldFileContent := ncr.Content
		ncr.Content = string(newFileContentBytes)
		return !strings.EqualFold(oldFileContent, ncr.Content), nil
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
