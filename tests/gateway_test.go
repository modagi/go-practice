package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"../gatewayserver"
)

func TestJsonRpc(t *testing.T) {
	var s gatewayserver.Server
	s.Init("../config.yml")
	go s.RunServer()

	protocolVersion := bytes.NewBufferString("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_protocolVersion\",\"params\":[]}")
	res, err := http.Post(fmt.Sprintf("http://127.0.0.1:%s", s.Config.Port), "application/json", protocolVersion)
	if err != nil {
		t.Error("http post failed")
	}
	if res.StatusCode != 200 {
		t.Error(fmt.Sprintf("status code: %d", res.StatusCode))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	var dat map[string]interface{}
	json.Unmarshal(buf.Bytes(), &dat)
	if dat["result"] != "63" {
		t.Error(fmt.Sprintf("protocolVersion result: %s", dat["result"]))
	}

}
