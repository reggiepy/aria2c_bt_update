package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gookit/goutil/byteutil"
	"github.com/reggiepy/aria2c_bt_updater/aria2c"
)

type Config struct {
	HttpProxy    string        `yaml:"http_proxy,inline" json:"http_proxy"`
	BtTrackerUrl []string      `yaml:"bt_tracker_url,inline" json:"bt_tracker_url"`
	Frequency    time.Duration `yaml:"frequency,inline" json:"frequency"`
}

type Server struct {
	jsonRpc      *aria2c.JsonRpc
	cfg          Config
	btTrackerMd5 string
}

func NewServer(jsonRpc *aria2c.JsonRpc, cfg Config) *Server {
	return &Server{
		jsonRpc: jsonRpc,
		cfg:     cfg,
	}
}

func (s *Server) CheckRpc() {
	for cnt := 1; ; cnt++ {
		if s.IsRunning() {
			break
		}
		zap.L().Info(fmt.Sprintf("服务未运行，重试 %d...", cnt))
		time.Sleep(time.Second)
	}
	return
}

// RefreshBtTracker updates the current BT Tracker if necessary.
func (s *Server) RefreshBtTracker() {
	for {
		result := s.jsonRpc.Post(aria2c.GetGlobalOption, nil, nil, aria2c.NewRpcOption())
		if result.Message == "failed" {
			zap.L().Info(fmt.Sprintf("Failed to get global options from RPC, retrying..."))
			time.Sleep(time.Second)
			continue
		}

		// Parse the result to get current BT Tracker
		data, ok := result.Result.(map[string]interface{})
		if !ok {
			zap.L().Info(fmt.Sprintf("Failed to parse result into map[string]interface{}"))
			continue
		}

		btTracker, ok := data["bt-tracker"].(string)
		if !ok {
			zap.L().Info(fmt.Sprintf("BT tracker is not a string"))
			continue
		}

		btTrackerMd5 := string(byteutil.Md5(btTracker))
		zap.L().Info(fmt.Sprintf("Current BT tracker MD5: %v", btTrackerMd5))
		s.btTrackerMd5 = btTrackerMd5
		break
	}
}

// UpdateBtTrackerUrls fetches BT trackers and updates them.
func (s *Server) UpdateBtTrackerUrls() {
	for index, url := range s.cfg.BtTrackerUrl {
		btTracker, err := s.GetBtTracker(url, 60*time.Second)
		if err != nil {
			zap.L().Info(fmt.Sprintf("Failed to get BT Tracker from URL: %s", url))
			continue
		}
		zap.L().Info(fmt.Sprintf("BtTracker %d (%s): %s", index, url, btTracker))
		err = s.SetBTTracker(btTracker)
		if err != nil {
			zap.L().Error("Failed to set BT Tracker", zap.String("tracker", btTracker), zap.Error(err))
		} else {
			zap.L().Info("Updated BT Tracker successfully", zap.String("url", url))
		}
	}
}

// SetBTTracker updates the BT tracker in the aria2c RPC service.
func (s *Server) SetBTTracker(btTracker string) error {
	// Calculate new BT tracker MD5
	newBtTrackerMd5 := string(byteutil.Md5(btTracker))

	// If MD5 is unchanged, skip update
	if newBtTrackerMd5 == s.btTrackerMd5 {
		return fmt.Errorf("BT Tracker MD5 is unchanged")
	}

	// Prepare the options for updating the tracker
	options := aria2c.Options{"bt-tracker": btTracker}

	// Call the RPC to change the global option
	result := s.jsonRpc.Post(aria2c.ChangeGlobalOption, nil, options, aria2c.NewRpcOption())

	// Marshal the result for logging
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}

	// Log the response
	zap.L().Info("Change BT Tracker", zap.String("response", string(jsonResult)))

	// Check the response and update the tracker if successful
	if result.Message == "success" {
		fmt.Sprintf("Updated BT Tracker successfully: %s", btTracker)
		s.btTrackerMd5 = newBtTrackerMd5
		return nil
	}

	// Log error if update fails
	zap.L().Error("Failed to update BT Tracker",
		zap.String("message", result.Message),
		zap.String("btTracker", btTracker))

	return fmt.Errorf("failed to update BT Tracker: %s", result.Message)
}

func (s *Server) IsRunning() bool {
	result := s.jsonRpc.Post(aria2c.GetVersion, nil, nil, aria2c.NewRpcOption())
	if result.Message == "failed" {
		zap.L().Error("RPC service is not running")
		return false
	}
	data, ok := result.Result.(map[string]interface{})
	if !ok || data["version"] == "" {
		zap.L().Error("Failed to retrieve RPC version")
		return false
	}
	return true
}

func (s *Server) Run() {
	s.CheckRpc()
	s.RefreshBtTracker()
	s.UpdateBtTrackerUrls()

	ticker := time.NewTicker(time.Second * s.cfg.Frequency)
	for {
		<-ticker.C
		s.UpdateBtTrackerUrls()
	}
}

// GetBtTracker retrieves the BT tracker list from a given URL.
func (s *Server) GetBtTracker(btTrackerUrl string, timeout time.Duration) (string, error) {
	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Set up proxy if configured
	proxy := http.ProxyFromEnvironment
	if s.cfg.HttpProxy != "" {
		fixedURL, err := url.Parse(s.cfg.HttpProxy)
		if err == nil {
			proxy = http.ProxyURL(fixedURL)
		}
	}

	// Reuse HTTP client to make the request
	req, err := http.NewRequestWithContext(ctx, "GET", btTrackerUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Create HTTP client with configured proxy
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// Read the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Return the body as a string
	return string(body), nil
}
