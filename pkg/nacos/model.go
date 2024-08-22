package nacos

import "time"

type Time struct {
	time.Time
}

func (n *Time) Show() string {
	return n.Time.Format("2006-01-02 15:04:05")
}
func (n *Time) MarshalJSON() ([]byte, error) {
	return []byte(n.Time.Format("2006-01-02T15:04:05.000+08:00")), nil
}
func (n *Time) UnmarshalJSON(bytes []byte) (err error) {
	n.Time, err = time.Parse(`"2006-01-02T15:04:05.000+08:00"`, string(bytes))
	return err
}

type StateResponse struct {
	AuthEnabled    string `json:"auth_enabled"`
	StandaloneMode string `json:"standalone_mode"`
	Version        string `json:"version"`
}

type LoginResponse struct {
	Token string `json:"accessToken"`
}

type NamespacesResponse struct {
	Code int              `json:"code"`
	Data []NamespacesItem `json:"data"`
}
type NamespacesItem struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	NamespaceDesc     string `json:"namespaceDesc"`
	Quota             int    `json:"quota"`
	ConfigCount       int    `json:"configCount"`
	Type              int    `json:"type"`
}

type ConfigsResponse struct {
	PageNumber int           `json:"pageNumber"`
	TotalCount int           `json:"totalCount"`
	PageItems  []ConfigsItem `json:"pageItems"`
}
type ConfigsItem struct {
	Content string `json:"content"`
	DataId  string `json:"dataId"`
	Group   string `json:"group"`
	AppName string `json:"appName"`
	Tenant  string `json:"tenant"`
	Id      string `json:"id"`
}

type ConfigResponse struct {
	Id               string `json:"id"`
	DataId           string `json:"dataId"`
	Group            string `json:"group"`
	Content          string `json:"content"`
	Md5              string `json:"md5"`
	EncryptedDataKey string `json:"encryptedDataKey"`
	Tenant           string `json:"tenant"`
	AppName          string `json:"appName"`
	CreateTime       int64  `json:"createTime"`
	ModifyTime       int64  `json:"modifyTime"`
	CreateUser       string `json:"createUser"`
	CreateIp         string `json:"createIp"`
	Desc             string `json:"desc"`
	Use              string `json:"use"`
	Effect           string `json:"effect"`
	Schema           string `json:"schema"`
	ConfigTags       string `json:"configTags"`
}

type ConfigDeleteResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    bool   `json:"data"`
}

type ConfigListenerModel struct {
	CollectStatus           int               `json:"collectStatus"`
	LisentersGroupkeyStatus map[string]string `json:"lisentersGroupkeyStatus"`
}

type HistoriesResponse struct {
	PageNumber int             `json:"pageNumber"`
	TotalCount int             `json:"totalCount"`
	PageItems  []HistoriesItem `json:"pageItems"`
}
type HistoriesItem struct {
	Id               string `json:"id"`
	LastId           int    `json:"lastId"`
	DataId           string `json:"dataId"`
	Group            string `json:"group"`
	Tenant           string `json:"tenant"`
	AppName          string `json:"appName"`
	Md5              string `json:"md5"`
	Content          string `json:"content"`
	SrcIp            string `json:"srcIp"`
	SrcUser          string `json:"srcUser"`
	OpType           string `json:"opType"`
	CreatedTime      Time   `json:"createdTime"`
	LastModifiedTime Time   `json:"lastModifiedTime"`
	EncryptedDataKey string `json:"encryptedDataKey"`
}

type ConfigCloneItem struct {
	CfgId  string `json:"cfgId"`
	DataId string `json:"dataId"`
	Group  string `json:"group"`
}

type ConfigCloneResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    ConfigCloneInfo `json:"data"`
}
type ConfigCloneInfo struct {
	SuccCount int               `json:"succCount"`
	SkipCount int               `json:"skipCount"`
	FailData  []ConfigCloneItem `json:"failData"`
	SkipData  []ConfigCloneItem `json:"skipData"`
}

type ServicesResponse struct {
	Count       int            `json:"count"`
	ServiceList []ServicesItem `json:"serviceList"`
}
type ServicesItem struct {
	Name                 string `json:"name"`
	GroupName            string `json:"groupName"`
	ClusterCount         int    `json:"clusterCount"`
	IpCount              int    `json:"ipCount"`
	HealthyInstanceCount int    `json:"healthyInstanceCount"`
	TriggerFlag          string `json:"triggerFlag"`
}

type ServiceResponse struct {
	Service  ServiceItem   `json:"service"`
	Clusters []ClusterItem `json:"clusters"`
}
type ServiceItem struct {
	Name             string            `json:"name"`
	GroupName        string            `json:"groupName"`
	ProtectThreshold float32           `json:"protectThreshold"`
	Selector         map[string]string `json:"selector"`
	Metadata         map[string]string `json:"metadata"`
}
type ClusterItem struct {
	ServiceName      string            `json:"serviceName"`
	Name             string            `json:"name"`
	HealthChecker    map[string]any    `json:"healthChecker"`
	DefaultPort      uint16            `json:"defaultPort"`
	DefaultCheckPort uint16            `json:"defaultCheckPort"`
	UseIPPort4Check  bool              `json:"useIPPort4Check"`
	Metadata         map[string]string `json:"metadata"`
}

type InstancesResponse struct {
	Count int             `json:"count"`
	List  []InstancesItem `json:"list"`
}
type InstancesItem struct {
	InstanceId                string            `json:"instanceId"`
	Ip                        string            `json:"ip"`
	Port                      uint16            `json:"port"`
	Weight                    float32           `json:"weight"`
	Healthy                   bool              `json:"healthy"`
	Enabled                   bool              `json:"enabled"`
	Ephemeral                 bool              `json:"ephemeral"`
	ClusterName               string            `json:"clusterName"`
	ServiceName               string            `json:"serviceName"`
	Metadata                  map[string]string `json:"metadata"`
	InstanceHeartBeatInterval int32             `json:"instanceHeartBeatInterval"`
	InstanceHeartBeatTimeOut  int32             `json:"instanceHeartBeatTimeOut"`
	IpDeleteTimeout           int32             `json:"ipDeleteTimeout"`
}

type SubscribersResponse struct {
	Count       int               `json:"count"`
	Subscribers []SubscribersItem `json:"subscribers"`
}
type SubscribersItem struct {
	AddrStr     string `json:"addrStr"`
	Agent       string `json:"agent"`
	App         string `json:"app"`
	Ip          string `json:"ip"`
	Port        uint16 `json:"port"`
	NamespaceId string `json:"namespaceId"`
	ServiceName string `json:"serviceName"`
	Cluster     string `json:"cluster"`
}
