package util

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type EditContentCallback func(ok bool, newContentYaml string, err error)
type EditStructCallback[T any] func(ok bool, newContent T, err error)

func EditContent(command string, fileName string, content string, fn EditContentCallback) tea.Cmd {
	filePath := filepath.Join(os.TempDir(), "nacos-tui", fileName)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		fn(false, "", errors.Join(errors.New("mkdir error"), err))
	}
	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		fn(false, "", errors.Join(errors.New("write file error"), err))
	}
	commandAndArgs := strings.Split(command, " ")
	args := commandAndArgs[1:]
	args = append(args, fileName)
	return tea.ExecProcess(exec.Command(commandAndArgs[0], args...), func(err error) tea.Msg {
		defer func(name string) {
			rErr := os.Remove(name)
			if rErr != nil {
				fn(false, "", errors.Join(errors.New("remove tmp file error"), rErr))
			}
		}(fileName)
		if err != nil {
			fn(false, "", errors.Join(errors.New("exec "+commandAndArgs[0]+" error"), err))
			return nil
		}
		newContentBytes, err := os.ReadFile(fileName)
		if err != nil {
			fn(false, "", errors.Join(errors.New("read tmp file error"), err))
		}
		newContent := string(newContentBytes)
		fn(!strings.EqualFold(content, newContent), newContent, nil)
		return base.RefreshScreenMsg
	})
}

func EditStruct[T any](command string, fileName string, content T, fn EditStructCallback[T]) tea.Cmd {
	var newContent T
	contentBytes, err := yaml.Marshal(content)
	if err != nil {
		fn(false, newContent, err)
	}
	return EditContent(command, fileName, string(contentBytes), func(ok bool, newContentYaml string, err error) {
		if err != nil {
			fn(false, newContent, err)
		}
		if ok {
			var err = yaml.Unmarshal([]byte(newContentYaml), &newContent)
			if err != nil {
				fn(false, newContent, err)
			}
			fn(ok, newContent, nil)
		}
		fn(ok, newContent, err)
	})
}
