package main

import (
	"fmt"
	"holders/conf"
	"holders/scanner"
	api "holders/service"
)

func main() {
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
