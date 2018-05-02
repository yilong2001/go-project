package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	Max_Request_Num int    = 50
	Default_PreStep int64  = 100
	IdGen_DB_Name   string = "web_idgens"
)

// 工具函数集

// 返回配置中的每次最大申请ID值
func getPreStep() int64 {
	return Default_PreStep
}

// 返回当前的运行目录
func runPath() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}

	return filepath.Dir(path)
}
