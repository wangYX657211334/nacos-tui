package nacos

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func createId() string {
	// 设置随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := make([]rune, 8)
	for i := range id {
		id[i] = rune(r.Intn(26) + 'a')
	}
	return string(id)
}

var id = createId()

// http客户端通用配置相关
type LoggerRoundTripper struct {
	BaseRoundTripper http.RoundTripper
	db               *sql.DB
	loginParam       func() (url string, username string, password string, namespace string, err error)
}

func (l *LoggerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var requestDupm, responseDupm string
	errMessage := errors.New("")
	baseUrl, username, password, namespace, err := l.loginParam()
	if err != nil {
		errMessage = errors.Join(errMessage, err)
	}
	dumpReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		errMessage = errors.Join(errMessage, err)
	} else {
		requestDupm = string(dumpReq)
	}

	res, err := l.BaseRoundTripper.RoundTrip(req)

	if err != nil {
		errMessage = errors.Join(errMessage, err)
	} else {
		dumpRes, err := httputil.DumpResponse(res, true)
		if err != nil {
			errMessage = errors.Join(errMessage, err)
		} else {
			responseDupm = string(dumpRes)
		}
	}
	_, err = l.db.Exec(`
INSERT INTO audit(session_id, base_url, username, password, namespace, request_dump, response_dump, error_message, time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, baseUrl, username, password, namespace, requestDupm, responseDupm, errMessage.Error(), time.Now().UnixMilli())
	if err != nil {
		log.Printf("save audit error: %s \n", err)
	}
	return res, err
}

type Api interface {
	ResetInitialization()
	GetNamespaces() (*NamespacesResponse, error)
	GetConfigs(dataId string, group string, pageNo int, pageSize int) (*ConfigsResponse, error)
	GetConfig(dataId string, group string) (*ConfigResponse, error)
	UpdateConfig(data *ConfigResponse) (bool, error)
	CloneConfigs(targetNamespace string, policy string, body []ConfigCloneItem) (*ConfigCloneResponse, error)
	DeleteConfig(ids ...string) (*ConfigDeleteResponse, error)
	GetConfigListener(dataId string, group string) (*ConfigListenerModel, error)
	GetConfigHistories(dataId string, group string, pageNo int, pageSize int) (*HistoriesResponse, error)
	GetServices(dataId string, group string, pageNo int, pageSize int) (*ServicesResponse, error)
	GetService(dataId string, group string) (*ServiceResponse, error)
	GetInstances(dataId string, group string, clusterName string) (*InstancesResponse, error)
	GetSubscribers(dataId string, group string, pageNo int, pageSize int) (*SubscribersResponse, error)
}

type api struct {
	url         string
	username    string
	password    string
	loginParam  func() (url string, username string, password string, namespace string, err error)
	token       string
	namespace   string
	initialized bool
	db          *sql.DB
	httpClient  *http.Client
}

func NewApi(loginParam func() (url string, username string, password string, namespace string, err error), db *sql.DB) Api {
	return &api{
		loginParam:  loginParam,
		initialized: false,
		db:          db,
		httpClient: &http.Client{
			Transport: &LoggerRoundTripper{
				BaseRoundTripper: http.DefaultTransport,
				db:               db,
				loginParam:       loginParam,
			},
			Timeout: 2 * time.Second,
		},
	}
}

func (n *api) ResetInitialization() {
	n.initialized = false
}

func (n *api) state() (*StateResponse, error) {
	return callRemote[StateResponse](n, func() (*http.Response, error) { return n.httpClient.Get(n.url + "/v1/console/server/state") })
}

const loginUrl = "/v1/auth/login"

func (n *api) refresh() (err error) {
	n.url, n.username, n.password, n.namespace, err = n.loginParam()
	if err != nil {
		return err
	}
	state, err := n.state()
	if err != nil {
		return err
	}
	if state.AuthEnabled != "true" {
		return nil
	}
	loginRes, err := callRemote[LoginResponse](n, func() (*http.Response, error) {
		return n.httpClient.PostForm(n.getUrl("/v1/auth/login"),
			url.Values{"username": {n.username}, "password": {n.password}})
	})
	if err != nil {
		return err
	}
	n.token = loginRes.Token
	return nil
}

func (n *api) GetNamespaces() (*NamespacesResponse, error) {
	return requiredLoginCallRemote[NamespacesResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/console/namespaces?accessToken=%s", n.token))
	})
}

func (n *api) GetConfigs(dataId string, group string, pageNo int, pageSize int) (*ConfigsResponse, error) {
	if len(dataId) > 0 {
		dataId = "*" + dataId + "*"
	}
	if len(group) > 0 {
		group = "*" + group + "*"
	}
	return requiredLoginCallRemote[ConfigsResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/cs/configs?dataId=%s&group=%s&search=blur&appName=&config_tags=&pageNo=%d&pageSize=%d&tenant=%s&accessToken=%s",
				dataId, group, pageNo, pageSize, n.namespace, n.token))
	})
}

func (n *api) GetConfig(dataId string, group string) (*ConfigResponse, error) {
	return requiredLoginCallRemote[ConfigResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/cs/configs?show=all&dataId=%s&group=%s&tenant=%s&accessToken=%s",
				dataId, group, n.namespace, n.token))
	})
}
func (n *api) DeleteConfig(ids ...string) (*ConfigDeleteResponse, error) {
	return requiredLoginCallRemote[ConfigDeleteResponse](n, func() (*http.Response, error) {
		request, err := http.NewRequest("DELETE",
			n.getUrl("/v1/cs/configs?delType=ids&ids=%s&tenant=%s&accessToken=%s", strings.Join(ids, ","),
				n.namespace, n.token), nil)
		if err != nil {
			return nil, err
		}
		return n.httpClient.Do(request)
	})
}

func (n *api) UpdateConfig(configRes *ConfigResponse) (bool, error) {
	res, err := n.httpClient.PostForm(n.url+fmt.Sprintf("/v1/cs/configs?accessToken=%s", n.token), url.Values{
		"id":               {configRes.Id},
		"dataId":           {configRes.DataId},
		"group":            {configRes.Group},
		"content":          {configRes.Content},
		"md5":              {configRes.Md5},
		"encryptedDataKey": {configRes.EncryptedDataKey},
		"tenant":           {configRes.Tenant},
		"appName":          {configRes.AppName},
		"createTime":       {strconv.FormatInt(configRes.CreateTime, 10)},
		"modifyTime":       {strconv.FormatInt(configRes.ModifyTime, 10)},
		"createUser":       {configRes.CreateUser},
		"createIp":         {configRes.CreateIp},
		"desc":             {configRes.Desc},
		"use":              {configRes.Use},
		"effect":           {configRes.Effect},
		"schema":           {configRes.Schema},
		"configTags":       {configRes.ConfigTags},
	})
	if err != nil {
		return false, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return false, errors.New("nacos request error")
	}
	return true, nil
}

func (n *api) GetConfigListener(dataId string, group string) (*ConfigListenerModel, error) {
	return requiredLoginCallRemote[ConfigListenerModel](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/cs/configs/listener?dataId=%s&group=%s&tenant=%s&accessToken=%s",
				dataId, group, n.namespace, n.token))
	})
}

func (n *api) GetConfigHistories(dataId string, group string, pageNo int, pageSize int) (*HistoriesResponse, error) {
	return requiredLoginCallRemote[HistoriesResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/cs/history?search=accurate&dataId=%s&group=%s&pageNo=%d&pageSize=%d&tenant=%s&accessToken=%s",
				dataId, group, pageNo, pageSize, n.namespace, n.token))
	})
}

func (n *api) CloneConfigs(targetNamespace string, policy string, body []ConfigCloneItem) (*ConfigCloneResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Join(errors.New("nacos request body marshal error"), err)
	}
	return requiredLoginCallRemote[ConfigCloneResponse](n, func() (*http.Response, error) {
		return n.httpClient.Post(
			n.getUrl("/v1/cs/configs?clone=true&tenant=%s&policy=%s&accessToken=%s",
				targetNamespace, policy, n.token), "application/json", bytes.NewReader(bodyBytes))
	})
}

func (n *api) GetServices(dataId string, group string, pageNo int, pageSize int) (*ServicesResponse, error) {
	return requiredLoginCallRemote[ServicesResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/ns/catalog/services?serviceNameParam=%s&groupNameParam=%s&hasIpCount=true&withInstances=false&pageNo=%d&pageSize=%d&namespaceId=%s&accessToken=%s",
				dataId, group, pageNo, pageSize, n.namespace, n.token))
	})
}

func (n *api) GetService(dataId string, group string) (*ServiceResponse, error) {
	return requiredLoginCallRemote[ServiceResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/ns/catalog/service?serviceName=%s&groupName=%s&namespaceId=%s&accessToken=%s",
				dataId, group, n.namespace, n.token))
	})
}

func (n *api) GetInstances(dataId string, group string, clusterName string) (*InstancesResponse, error) {
	return requiredLoginCallRemote[InstancesResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/ns/catalog/instances?serviceName=%s&groupName=%s&clusterName=%s&pageNo=1&pageSize=1000&namespaceId=%s&accessToken=%s",
				dataId, group, clusterName, n.namespace, n.token))
	})
}

func (n *api) GetSubscribers(dataId string, group string, pageNo int, pageSize int) (*SubscribersResponse, error) {
	return requiredLoginCallRemote[SubscribersResponse](n, func() (*http.Response, error) {
		return n.httpClient.Get(
			n.getUrl("/v1/ns/service/subscribers?serviceName=%s&groupName=%s&pageNo=%d&pageSize=%d&namespaceId=%s&accessToken=%s",
				dataId, group, pageNo, pageSize, n.namespace, n.token))
	})
}

func (n *api) getUrl(format string, a ...any) string {
	return n.url + fmt.Sprintf(format, a...)
}
func requiredLoginCallRemote[T any](n *api, call func() (*http.Response, error)) (*T, error) {
	if !n.initialized {
		if err := n.refresh(); err != nil {
			return nil, err
		}
		n.initialized = true
	}
	return callRemote[T](n, call)
}
func callRemote[T any](n *api, call func() (*http.Response, error)) (res *T, err error) {
	// 实际调用
	response, err := call()
	if err != nil {
		return nil, errors.Join(errors.New("nacos request error"), err)
	}
	if response.StatusCode == 403 {
		if strings.Contains(response.Request.RequestURI, loginUrl) {
			return nil, errors.Join(errors.New("nacos login error"), err)
		}
		// 重新登陆
		if err = n.refresh(); err != nil {
			return nil, err
		}
		// 重试一次
		if response, err = call(); err != nil {
			return nil, errors.Join(errors.New("nacos request error"), err)
		}
	}
	defer func(body io.ReadCloser) {
		if closeErr := body.Close(); closeErr != nil {
			err = errors.Join(errors.New("nacos response body close error"), closeErr, err)
			return
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Join(errors.New("nacos response body read error"), err)
	}
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.Join(errors.New("nacos response body json unmarshal error"), err)
	}
	return res, nil
}
