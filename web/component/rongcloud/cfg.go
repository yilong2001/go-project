package rongcloud

import (
	"encoding/json"
)

type RongCloudAppConfig struct {
	AppId  string
	AppKey string
}

func GetRongCloudAppConfig() *RongCloudAppConfig {
	return &RongCloudAppConfig{
		AppId:  "k51hidwq1qp0b",
		AppKey: "R1HLJyKQFc",
	}
}

type RongCloudResult struct {
	Code   int    `json:"code"`
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func GetRongCloudResult(res string) *RongCloudResult {
	result := &RongCloudResult{}

	err := json.Unmarshal([]byte(res), result)
	if err != nil {
		return nil
	}

	return result
}
