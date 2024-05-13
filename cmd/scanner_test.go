package main

import (
	"fmt"
	"holders/conf"
	"holders/db"
	"holders/scanner"
	api "holders/service"
	"testing"
)

func TestNumber(t *testing.T) {
	err := db.WriteNumber("btc-testNet","2810930")
	if err != nil {
		fmt.Println(err)
	}
}

func TestScanner(t *testing.T) {
	client, err := scanner.NewClient(conf.NodeUrl, "btc-testNet")
	if err != nil {
		fmt.Println(err)
	}

	//扫描日志
	go client.FilterLogs()

	//解析日志
	go client.ResolveLogs()

	service := api.NewGinService()
	service.Run(":8085")
}

func TestService(t *testing.T) {
	service := api.NewGinService()
	service.Run(":8085")
}


func TestServer(t *testing.T) {
	service := api.NewGinService()
	service.Run(":8085")
}