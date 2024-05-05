package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"holders/conf"
	"holders/db"
	"holders/jsonrpc"
	"holders/models"
	"log"
	"testing"
	"time"
)

func TestBestBlockNumber(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://152.53.33.145:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}
	result, err := cli.BestBlockNumber()
	fmt.Println(result)
}

func TestCallContract(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}

	params := []string{"1"}

	param := jsonrpc.CallParam{
		KID:    "kfc3f6e0ec06038784e17876450ee6af8d32132dd5",
		Method: "$balanceOf",
		Params:   params,
	}
	result, err := cli.CallContract(param)
	fmt.Println(result, err)
}


func TestEventForBlock(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}
	param := jsonrpc.EventParam{
		Number: "11906938",
	}
	result, err := cli.GetEvents(param)
	fmt.Println(result, err)
}

func TestBockNumber(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}

	param := jsonrpc.BlockNumberParam{
		Number: "11906938",
	}
	result, err := cli.GetBlockNumber(param)
	if err != nil {
		log.Println(err)
	}
	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}

func TestTransaction(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}
	param := jsonrpc.TransactionParam{
		Hash: "af52c8aec1e6d5a47ad95c3d3393afb811e11f4f32be7aeebfbe517ffc5da13f",
	}
	result, err := cli.GetTransaction(param)
	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}

func TestScriptModel(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}
	param := jsonrpc.ScriptParam{
		KID: "kfcf99351536cc66b3aeb8b39be3e1e44e4fd78193",
	}
	result, err := cli.GetScriptModel(param)
	abi := result.(*jsonrpc.Script)
	fmt.Println(abi.Kip, err)
}

func TestTokenModel(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}
	param := jsonrpc.TokenParam{
		KID: "kfc2449ae72e89f90e11370502b6992af7c85aa1fa",
	}
	result, err := cli.GetTokenModel(param)
	token := result.(*jsonrpc.Token)
	fmt.Println(token, err)
}

func TestTokenUri(t *testing.T) {
	cli, err := jsonrpc.NewClient("http://127.0.0.1:7399/jrpc")
	if err != nil {
		fmt.Println(err)
	}
	param := jsonrpc.TokenUriParam{
		KID:     "kfc1715339bf254ee12fb03da6ba1099cd831e9d2b",
		TokenId: "1000",
	}
	result, err := cli.GetTokenUri(param)
	token := result.(string)
	fmt.Println(token, err)
}



func Test24H(t *testing.T) {
	nowUnix := time.Now()
	h, _ := time.ParseDuration("-1h")
	Unix24 := nowUnix.Add(24 * h).Unix()
	fmt.Println(Unix24)
}

func TestUUID(t *testing.T) {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())
	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}



func TestT20(t *testing.T) {
	t20 := models.Transfer20{
		Amount: 1000,
		Kid:    "kid",
		From:   "to1",
		To:     "to2",
	}
	db.Transaction20(t20)
}

func TestT721(t *testing.T) {
	transfer721 := models.Transfer721{
		Kid:     "kid7",
		From:    "to",
		To:      "to1",
		TokenId: "10000",
	}
	db.Transaction721(transfer721)
}

func TestToken(t *testing.T) {
	getTokenMeta("kfc2449ae72e89f90e11370502b6992af7c85aa1fa")
}

func getTokenMeta(kid string) {
	exits := db.GetTokenExits(kid)
	fmt.Println(exits)
	if !exits {
		//获取对应信息并保存
		rpc, err := jsonrpc.NewClient(conf.NodeUrl)
		if err != nil {
			return
		}
		param := jsonrpc.TokenParam{
			KID: kid,
		}

		token, err := rpc.GetTokenModel(param)
		if err != nil {
			fmt.Println(123)
			log.Println(err)
			return
		}

		t := token.(*jsonrpc.Token)

		t2 := models.Token{
			Kid:         kid,
			Name:        t.Name,
			Symbol:      t.Symbol,
			TotalSupply: t.TotalSupply,
		}

		err = db.Token(t2)
		if err != nil {
			return
		}
		db.PutTokenExits(kid)
	}
}

func TestHold(t *testing.T) {
	holds, err := db.FindWalletHold("2N7TYrDKNeZf4eVGXDVJyRKWaPdbx4qvCJj")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(holds)
}

func TestTokenIds(t *testing.T) {
	tokenIds, err := db.FindTokenIds("kfc1715339bf254ee12fb03da6ba1099cd831e9d2b", "wallet1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenIds)
}


func TestDist(t *testing.T) {
	dist, err := db.FindDist("kfc1715339bf254ee12fb03da6ba1099cd831e9d2b", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dist)
}
