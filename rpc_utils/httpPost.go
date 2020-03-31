package rpc_utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func HttpPostWResponse(reqBytes []byte) (Response, error) {
	req, err := http.NewRequest("POST", g_ethRpc.Url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return *new(Response), err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return *new(Response), err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return *new(Response), err
	}
	return response, nil
}
