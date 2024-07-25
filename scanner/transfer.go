package scanner

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"holders/conf"
	"holders/db"
	"holders/jsonrpc"
	"holders/models"
	"log"
	"strconv"
)

func transfer(e jsonrpc.Event) {
	if e.Name == "Transfer" {
		cli := jsonrpc.GetClient()
		param := jsonrpc.ScriptParam{
			KID: e.KID,
		}
		result, err := cli.GetScriptModel(param)
		if err != nil {
			log.Println(err)
			return
		}

		var checkKid string
		script := result.(*jsonrpc.Script)
		switch script.Kip {
		case "B20":
			sAmount := fmt.Sprint(e.Args["amount"])
			// 将字符串转换为float64
			amount, err := strconv.ParseFloat(sAmount, 64)
			if err != nil {
				fmt.Println("转换错误:", err)
				return
			}
			fmt.Println(e.Args)
			if e.Args["from"] == nil || e.Args["to"] == nil {
				return
			}
			//记录K20转账
			t20 := models.Transfer20{
				Kid: e.KID,
				From: e.Args["from"].(string),
				To: e.Args["to"].(string),
				Amount: amount,
			}
			err = db.Transaction20(t20)
			if err != nil {
				log.Println(err)
				return
			}
			checkKid = t20.Kid
		case "B721":
			//记录K721转账
			var t721 models.Transfer721
			err := mapstructure.Decode(e.Args, &t721)
			if err != nil {
				fmt.Println(err)
				break
			}
			t721.Kid = e.KID

			uri, err := getTokenUri(t721.Kid, fmt.Sprint(t721.TokenId))
			if err != nil {
				log.Println(err)
			}
			t721.Data = uri

			err = db.Transaction721(t721)
			if err != nil {
				log.Println(err)
				return
			}
			checkKid = t721.Kid
		}

		//获取代币信息
		if checkKid != "" {
			go getTokenMeta(checkKid)
		}

	}
}

// 获取该合约额外信息
func getTokenMeta(kid string) {
	exits := db.GetTokenExits(kid)
	if !exits {
		t2 := models.Token{
			Kid: kid,
		}

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
			log.Println(err)
		}

		if token != nil {
			t := token.(*jsonrpc.Token)
			t2.Name = t.Name
			t2.Symbol = t.Symbol
			t2.TotalSupply = t.TotalSupply
		} else {
			t2.Name = "Unknown"
			t2.Symbol = "Unknown"
			t2.TotalSupply = "Unknown"
		}

		err = db.Token(t2)

		if err != nil {
			log.Println(err)
			return
		}
		db.PutTokenExits(kid)
	}
}

func getTokenUri(kid, tokenId string) (string, error) {
	exits := db.GetTokenExits(kid)
	if !exits {
		db.PutTokenUriExits(kid, tokenId)
		//获取对应信息并保存
		rpc, err := jsonrpc.NewClient(conf.NodeUrl)
		if err != nil {
			return "Unknown", err
		}
		param := jsonrpc.TokenUriParam{
			KID:     kid,
			TokenId: tokenId,
		}

		uri, err := rpc.GetTokenUri(param)
		if err != nil {
			return "Unknown", err
		}
		s := uri.(*string)
		return *s, nil
	} else {
		return "Unknown", nil
	}
}
