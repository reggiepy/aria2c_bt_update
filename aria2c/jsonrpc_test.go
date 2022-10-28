package aria2c

import (
	"encoding/json"
	"testing"
	"time"
)

const (
	host = "localhost"
	port = 6800
)

var (
	jsonRpcOption = JsonRpcOption{}
)

func TestGetGlobalStat(t *testing.T) {
	json_rpc := NewJsonRpc(host, port, "", jsonRpcOption)
	result := json_rpc.Post(GetGlobalStat, nil, nil, NewRpcOption())
	jsonResult, _ := json.Marshal(result)
	if result.Message != "success" {
		t.Error("result: ", string(jsonResult))
	} else {
		t.Log("result: ", string(jsonResult))
	}
}

func TestUpdateBtTracker(t *testing.T) {
	json_rpc := NewJsonRpc(host, port, "", jsonRpcOption)

	url := "https://cdn.staticaly.com/gh/XIU2/TrackersListCollection/master/best_aria2.txt"
	rc, btTracker := json_rpc.GetBtTracker(url, time.Second*60)
	if rc != 0 {
		t.Error("Get BtTracker error: ", rc, btTracker)
	}
	t.Log("btTracker: ", btTracker)

	params := Params{}
	params = append(params, map[string]string{"bt-tracker": btTracker})
	result := json_rpc.Post(ChangeGlobalOption, params, nil, NewRpcOption())
	jsonResult, _ := json.Marshal(result)
	if result.Message != "success" {
		t.Error("result: ", string(jsonResult))
	} else {
		t.Log("result: ", string(jsonResult))
	}
}

func TestGetGlobalOption(t *testing.T) {
	json_rpc := NewJsonRpc(host, port, "", jsonRpcOption)
	result := json_rpc.Post(GetGlobalOption, Params{}, nil, NewRpcOption())
	jsonResult, _ := json.Marshal(result)
	if result.Message != "success" {
		t.Error("result: ", string(jsonResult))
	} else {
		t.Log("result: ", string(jsonResult))
	}
}
