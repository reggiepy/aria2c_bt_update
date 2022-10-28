package main

import (
	"aria2_cbt_tracker_updater/common"
	"aria2_cbt_tracker_updater/pkg/util"
	"aria2_cbt_tracker_updater/pkg/util/log"
	"aria2_cbt_tracker_updater/pkg/yaml"
	"encoding/json"
	"github.com/reggie/aria2c"
	"time"
)

type Server struct {
	rpc          *aria2c.JsonRpc
	config       *common.Config
	btTrackerMd5 string
}

func NewServer() *Server {
	config := &common.Config{}
	err := common.CheckConfig()
	if err != nil {
		panic(err)
	}
	yaml.InitConfigByViper(common.DefaultConfigPath, config)

	logConfig := config.Log
	log.InitLog(
		logConfig.LogWay,
		logConfig.LogFile,
		logConfig.LogLevel,
		logConfig.LogMaxDays,
		logConfig.DisableLogColor,
	)
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
	result := s.rpc.Post(aria2c.GetGlobalOption, nil, nil, aria2c.NewRpcOption())
	jsonResult, _ := json.Marshal(result)
	log.Debug("GetGlobalOption: %s", string(jsonResult))
	data := result.Result.(map[string]interface{})
	btTrackerMd5 := util.MD5(data["bt-tracker"].(string))
	log.Info("current md5: %v", btTrackerMd5)
	s.btTrackerMd5 = btTrackerMd5
}

func (s *Server) UpdateGlobalOption() {
	rc, btTracker := s.rpc.GetBtTracker(s.config.BtTrackerUrl, time.Second*60)
	if rc != 0 {
		log.Info("Failed to get BtTracker %v", btTracker)
	} else {
		newBtTrackerMd5 := util.MD5(btTracker)
		if newBtTrackerMd5 != s.btTrackerMd5 {
			options := aria2c.Options{"bt-tracker": btTracker}
			result := s.rpc.Post(aria2c.ChangeGlobalOption, nil, options, aria2c.NewRpcOption())
			jsonResult, _ := json.Marshal(result)
			log.Debug("result: %v", string(jsonResult))
			if result.Message == "success" {
				log.Info("update btTracker success: %s %s", newBtTrackerMd5, btTracker)
				s.btTrackerMd5 = newBtTrackerMd5
			}
		}
	}
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
func main() {
	server := NewServer()
	server.Run()
}
