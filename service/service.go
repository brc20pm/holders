package api

import (
	"holders/conf"
	"holders/jsonrpc"
)

func call(param jsonrpc.CallParam) (any, error) {
	cli, err := jsonrpc.NewClient(conf.NodeUrl)
	if err != nil {
		return nil, err
	}
	return cli.CallContract(param)
}

