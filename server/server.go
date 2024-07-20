package server

import (
	"encoding/json"
	"fmt"
	"time"

	"aria2c_bt_updater/common"
	"aria2c_bt_updater/pkg/util/log"
	util2 "aria2c_bt_updater/pkg/util/util"

	"github.com/reggie/aria2c"
)

type Server struct {
	rpc          *aria2c.JsonRpc
	config       *common.Config
	btTrackerMd5 string
}

func NewServer(config *common.Config) *Server {
	jsonByte, _ := json.Marshal(&config)
	log.Debug("src: %v\n", string(jsonByte))
	jsonRpcOption := aria2c.JsonRpcOption{
		ProxyUrl: config.HttpProxy,
	}
	jsonRpc := aria2c.NewJsonRpc(
		config.App.Host,
		config.App.Port,
		config.App.Token,
		jsonRpcOption,
	)

	return &Server{
		rpc:    jsonRpc,
		config: config,
	}
}

func (s *Server) Init() {
	for cnt := 1; ; cnt++ {
		if !s.IsRunning() {
			log.Error(fmt.Sprintf("服务未运行，重试 %d...", cnt))
			continue
		}
		result := s.rpc.Post(aria2c.GetGlobalOption, nil, nil, aria2c.NewRpcOption())
		if result.Message == "failed" {
			time.Sleep(time.Second)
			log.Error(fmt.Sprintf("请求rpc服务失败，重试 %d...", cnt))
			continue
		}

		jsonResult, err := json.Marshal(result)
		if err != nil {
			log.Error("JSON 编码失败: %v", err)
			continue
		}
		log.Info("GetGlobalOption: %s", string(jsonResult))

		data, ok := result.Result.(map[string]interface{})
		if !ok {
			log.Error("无法将结果转换为 map[string]interface{}")
			continue
		}

		btTracker, ok := data["bt-tracker"].(string)
		if !ok {
			log.Error("bt-tracker 不是 string 类型")
			continue
		}

		btTrackerMd5 := util2.MD5(btTracker)
		log.Info("current md5: %v", btTrackerMd5)
		s.btTrackerMd5 = btTrackerMd5
		break // 如果需要在成功获取后退出循环，使用 break
	}
}

func (s *Server) UpdateGlobalOption() {
	if !s.IsRunning() {
		log.Error(fmt.Sprintf("服务未运行"))
		return
	}
	rc, btTracker := s.rpc.GetBtTracker(s.config.BtTrackerUrl, time.Second*60)
	if rc != 0 {
		log.Info("Failed to get BtTracker: %v", btTracker)
		return
	}

	newBtTrackerMd5 := util2.MD5(btTracker)
	if newBtTrackerMd5 == s.btTrackerMd5 {
		log.Info("BtTracker MD5 is unchanged")
		return
	}

	options := aria2c.Options{"bt-tracker": btTracker}
	result := s.rpc.Post(aria2c.ChangeGlobalOption, nil, options, aria2c.NewRpcOption())

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Error("Failed to marshal result: %v", err)
		return
	}
	log.Debug("Result: %s", string(jsonResult))

	if result.Message == "success" {
		log.Info("Updated btTracker successfully: %s %s", newBtTrackerMd5, btTracker)
		s.btTrackerMd5 = newBtTrackerMd5
	} else {
		log.Error("Failed to update btTracker: %s", result.Message)
	}
}

func (s *Server) IsRunning() bool {
	result := s.rpc.Post(aria2c.GetVersion, nil, nil, aria2c.NewRpcOption())
	if result.Message == "failed" {
		return false
	}
	data, ok := result.Result.(map[string]interface{})
	if !ok {
		return false
	}
	if data["version"] == "" {
		return false
	}
	return true
}

func (s *Server) Run() {
	s.Init()
	s.UpdateGlobalOption()

	ticker := time.NewTicker(time.Second * s.config.Frequency)
	for {
		<-ticker.C
		s.UpdateGlobalOption()
	}
}
