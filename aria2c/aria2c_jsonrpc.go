package aria2c

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	Header  = map[string]string{"Content-Type": "application/json"}
	Version = "2.0"
)

const (
	GetGlobalStat      = "aria2.getGlobalStat"
	ChangeGlobalOption = "aria2.changeGlobalOption"
	GetGlobalOption    = "aria2.getGlobalOption"
	GetVersion         = "aria2.getVersion"
)

type JsonRpc struct {
	Host  string
	Port  int
	Token string
	Id    string

	Url string

	ctx          context.Context
	cancel       context.CancelFunc
	detailClient *http.Client
}

type JsonRpcOption struct {
	ProxyUrl string `json:"ProxyUrl"`
}

func NewJsonRpc(host string, port int, token string, options JsonRpcOption) *JsonRpc {
	idString := uuid.NewString()
	ctx, cancel := context.WithCancel(context.Background())
	proxy := http.ProxyFromEnvironment
	if options.ProxyUrl != "" {
		fixedURL, err := url.Parse(options.ProxyUrl)
		if err == nil {
			proxy = http.ProxyURL(fixedURL)
		}
	}
	return &JsonRpc{
		Host:   host,
		Port:   port,
		Token:  token,
		Id:     idString,
		Url:    fmt.Sprintf("http://%s:%d/jsonrpc", host, port),
		ctx:    ctx,
		cancel: cancel,
		detailClient: &http.Client{
			Transport: &http.Transport{
				Proxy: proxy,
			},
		},
	}
}

type (
	Params  []interface{}
	Options map[string]interface{}
)

type Data struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type RpcOption struct {
	Timeout   time.Duration
	OnSuccess func()
	OnFailed  func()
}

func NewRpcOption() RpcOption {
	return RpcOption{
		Timeout: time.Second * 60,
	}
}

type Result struct {
	Id      string      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`

	// Additional Information
	Message string `json:"message"`
}

func (c *JsonRpc) Post(action string, params Params, options Options, rpcOption RpcOption) (result *Result) {
	result = &Result{}
	timeout := time.Second * 60
	if rpcOption.Timeout != 0 {
		timeout = rpcOption.Timeout
	}
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	defer func() {
		if result.Message == "" {
			result.Message = "success"
			if rpcOption.OnSuccess != nil {
				rpcOption.OnSuccess()
			}
		} else {
			result.Message = "failed"
			if rpcOption.OnFailed != nil {
				rpcOption.OnFailed()
			}
		}
	}()

	data := c.GenData(action, params, options)
	body, err := json.Marshal(data)
	if err != nil {
		result.Message = fmt.Sprintf("json.Marshaler error: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Url, strings.NewReader(string(body)))
	if err != nil {
		result.Message = fmt.Sprintf("handle req error: %v", err)
		return
	}

	for k, v := range Header {
		req.Header.Set(k, v)
	}

	response, err := c.detailClient.Do(req)
	if err != nil {
		result.Message = fmt.Sprintf("handle http error: %v", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		result.Message = fmt.Sprintf("request status not is 200 %v", response.StatusCode)
		return
	}

	resp_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		result.Message = fmt.Sprintf("parse response error: %v", err)
		return
	}

	json.Unmarshal(resp_body, &result)
	return
}

func (c *JsonRpc) GenData(action string, params Params, options Options) *Data {
	if params == nil {
		params = Params{}
	}
	if c.Token != "" {
		params = append(params, fmt.Sprintf("token:%s", c.Token))
	}
	if options != nil {
		params = append(params, options)
	}
	return &Data{
		Jsonrpc: Version,
		Id:      c.Id,
		Method:  action,
		Params:  params,
	}
}
