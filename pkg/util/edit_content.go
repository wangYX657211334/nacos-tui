package util

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func EditContentBySystemEditor(fileName string, content string) (_ bool, _ string, err error) {
	filePath := filepath.Join(os.TempDir(), "nacos-tui", fileName)
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return false, "", errors.Join(errors.New("mkdir error"), err)
	}
	if err = os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return false, "", errors.Join(errors.New("write file error"), err)
	}
	defer func(name string) {
		rErr := os.Remove(name)
		if rErr != nil {
			err = errors.Join(errors.New("remove tmp file error"), rErr, err)
		}
	}(fileName)
	cmd := exec.Command("vim", fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return false, "", errors.Join(errors.New("edit file error"), err)
	}
	newContentBytes, err := os.ReadFile(fileName)
	if err != nil {
		return false, "", errors.Join(errors.New("read tmp file error"), err)
	}
	newContent := string(newContentBytes)
	return !strings.EqualFold(content, newContent), newContent, nil
}

func EditStructBySystemEditor[T any](fileName string, content T) (_ bool, _ T, err error) {
	var newContent T
	contentBytes, err := yaml.Marshal(content)
	if err != nil {
		return false, newContent, err
	}
	ok, newContentYaml, err := EditContentBySystemEditor(fileName, string(contentBytes))
	if err != nil {
		return false, newContent, err
	}
	if ok {
		var err = yaml.Unmarshal([]byte(newContentYaml), &newContent)
		if err != nil {
			return false, newContent, err
		}
		return ok, newContent, nil
	}
	return ok, newContent, err
}
