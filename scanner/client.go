package scanner

import (
	"fmt"
	"holders/db"
	"holders/jsonrpc"
	"log"
	"time"
)

type rpc struct {
	client *jsonrpc.Client
	chain  string
	eChan  chan []jsonrpc.Event
}

func NewClient(url, chain string) (*rpc, error) {
	cli, err := jsonrpc.NewClient(url)
	if err != nil {
		return nil, err
	}
	return &rpc{client: cli, chain: chain, eChan: make(chan []jsonrpc.Event)}, nil
}

func (r *rpc) FilterLogs() {
	scanNumber := int64(0)

	for {
		number, err := r.client.BestBlockNumber()
		if err != nil {
			log.Println(err)
			continue
		}
		fistNumber := db.FistNumber(r.chain)
		lastNumber := number.(int64)
		localNumber := int64(fistNumber)

		if localNumber == 0 {
			//协议运行区块 - 1
			localNumber = 2588380
		}

		////模拟需要同步的区块
		//lastNumber = 12228573

		if localNumber < lastNumber {
			switch {
			case lastNumber-localNumber == 1:
				//如果落后一个,则用最新的区块号去获取日志和区块时间戳相关信息
				scanNumber = lastNumber
			case localNumber == 0:
				//如果本地同步的区块号为0,则用最新的区块号去获取日志和区块时间戳相关信息
				scanNumber = lastNumber
			case lastNumber-localNumber > 1:
				scanNumber = localNumber + 1
			}

			fmt.Println("indexed number is ",scanNumber)

			param := jsonrpc.EventParam{
				Number: fmt.Sprint(scanNumber),
			}
			events, err := r.client.GetEvents(param)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if events != nil {
				eventList := events.([]jsonrpc.Event)
				if eventList != nil {
					r.eChan <- eventList
				}
			}

			err = db.WriteNumber(r.chain, param.Number)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			time.Sleep(1 * time.Minute)
		}
	}
}

func (r *rpc) ResolveLogs() {
	for events := range r.eChan {
		for _, e := range events {
			transfer(e)
		}
	}
}
