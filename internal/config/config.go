package config

import (
	"database/sql"
	"errors"
	"strconv"
)

type Api interface {
	GetNacosContexts() ([]NacosContext, error)
	GetNacosContext() (NacosContext, error)
	AddNacosContext(context NacosContext) error
	UpdateNacosContext(context NacosContext) error
	SetActiveNacosContext(name string) error
	SetNacosContextNamespace(namespace string, namespaceName string) error
	GetProperty(key string, defaultValue string) (string, error)
	GetIntProperty(key string, defaultValue string) (int, error)
	GetBoolProperty(key string, defaultValue string) (bool, error)
	SetProperty(string, string) error
}

type configApi struct {
	db *sql.DB
}

func NewApi(db *sql.DB) Api {
	return &configApi{db: db}
}

type NacosContext struct {
	Name             string
	Url              string
	User             string
	Password         string
	UseNamespace     string
	UseNamespaceName string
}

func (c *configApi) GetNacosContexts() (_ []NacosContext, err error) {
	rows, err := c.db.Query(`select name, url, username, password, namespace, namespace_name from nacos_context`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(rows)
	var contexts []NacosContext
	for rows.Next() {
		var name, url, username, password, namespace, namespaceName string
		if err = rows.Scan(&name, &url, &username, &password, &namespace, &namespaceName); err != nil {
			return nil, err
		}
		contexts = append(contexts, NacosContext{
			Name:             name,
			Url:              url,
			User:             username,
			Password:         password,
			UseNamespace:     namespace,
			UseNamespaceName: namespaceName,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return contexts, nil
}

func (c *configApi) GetProperty(key string, defaultValue string) (_ string, err error) {
	rows, err := c.db.Query(`select value from system_config where key = ?`, key)
	if err != nil {
		return defaultValue, err
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(rows)
	if rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}
func (c *configApi) GetIntProperty(key string, defaultValue string) (int, error) {
	stringValue, err := c.GetProperty(key, defaultValue)
	if err != nil {
		return 0, err
	}
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
func (c *configApi) GetBoolProperty(key string, defaultValue string) (bool, error) {
	stringValue, err := c.GetProperty(key, defaultValue)
	if err != nil {
		return false, err
	}
	boolValue, err := strconv.ParseBool(stringValue)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}
func (c *configApi) SetProperty(key string, value string) (err error) {
	dbValue, err := c.GetProperty(key, "nil")
	if dbValue != "nil" {
		// 修改
		_, err := c.db.Exec(`update system_config set value = ? where key = ?`, value, key)
		if err != nil {
			return err
		}
	} else {
		// 新增
		_, err := c.db.Exec(`insert into system_config(key, value) values (?, ?)`, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *configApi) GetNacosContext() (NacosContext, error) {
	contexts, err := c.GetNacosContexts()
	if err != nil {
		return NacosContext{}, err
	}
	activeContextName, err := c.GetProperty("active.nacos.context", "")
	if err != nil {
		return NacosContext{}, err
	}
	for _, context := range contexts {
		if context.Name == activeContextName {
			return context, nil
		}
	}
	return NacosContext{}, errors.New("not found active nacos context")
}

func (c *configApi) UpdateNacosContext(context NacosContext) error {
	_, err := c.db.Exec(`update nacos_context 
set url = ?, username = ?, password = ?, namespace = ?, namespace_name = ? 
where name = ?`, context.Url, context.User, context.Password, context.UseNamespace, context.UseNamespaceName, context.Name)
	if err != nil {
		return err
	}
	return nil
}

func (c *configApi) AddNacosContext(context NacosContext) error {
	_, err := c.db.Exec(`insert into nacos_context(name, url, username, password, namespace, namespace_name) values (?, ?, ?, ?, ?, ?)`,
		context.Name, context.Url, context.User, context.Password, context.UseNamespace, context.UseNamespaceName)
	if err != nil {
		return err
	}
	return nil
}

func (c *configApi) SetActiveNacosContext(name string) error {
	return c.SetProperty("active.nacos.context", name)
}
func (c *configApi) SetNacosContextNamespace(namespace string, namespaceName string) error {
	activeContext, err := c.GetNacosContext()
	if err != nil {
		return err
	}
	_, err = c.db.Exec(`update nacos_context set namespace = ?, namespace_name = ? where name = ?`,
		namespace, namespaceName, activeContext.Name)
	if err != nil {
		return err
	}
	return nil
}
