package jsonrpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"io"
	"net/http"
	"strings"
)

// JSONRPCRequest 定义JSON-RPC请求的结构体
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      string          `json:"id"`
}

// JSONRPCResponse 定义JSON-RPC响应的结构体
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  JSONResult    `json:"result"`
	Error   *JSONRPCError `json:"error"`
	ID      string        `json:"id"`
}

type JSONResult struct {
	Data interface{} `json:"data"`
}

// JSONRPCError 定义JSON-RPC错误的结构体
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	url string `json:"url"`
}

var rpcClient *Client

func NewClient(nodeUrl string) (*Client, error) {
	if nodeUrl == "" {
		return nil, errors.New("err: nodeUrl invalid")
	}
	if !strings.HasPrefix(nodeUrl, "http://") && !strings.HasPrefix(nodeUrl, "https://") {
		return nil, errors.New("err: only http or https requests are supported")
	}

	rpcClient = &Client{url: nodeUrl}

	return rpcClient, nil
}

func GetClient() *Client {
	return rpcClient
}

// 发送JSON-RPC请求的函数
func sendJSONRPCRequest(url string, request JSONRPCRequest) (*JSONRPCResponse, error) {
	// 将请求结构体编码为JSON
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// 创建一个HTTP客户端
	client := &http.Client{}
	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	// 设置Content-Type为application/json
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned HTTP status %s", resp.Status)
	}

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 将响应体解码为JSON-RPC响应结构体
	var response JSONRPCResponse

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}
	// 如果响应包含错误，返回错误
	if response.Error != nil {
		data, _ := json.Marshal(response.Error)
		return nil, errors.New(string(data))
	}

	return &response, nil
}

func (c *Client) CallContract(param CallParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("ord_call", pByte)
}

// 获取节点处理完成的最新区块号
func (c *Client) BestBlockNumber() (any, error) {
	return c.Call("bestBlockNumber", nil)
}

// 获取脚本模型
func (c *Client) GetScriptModel(param ScriptParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("getScriptModel", pByte)
}

// 获取代币模型
func (c *Client) GetTokenModel(param TokenParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("getTokenModel", pByte)
}

// 获取代币模型
func (c *Client) GetTokenUri(param TokenUriParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("getTokenUri", pByte)
}

// 获取指定合约地址在区块当中的事件记录
func (c *Client) GetEvents(param EventParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("getEvents", pByte)
}

func (c *Client) GetBlockNumber(param BlockNumberParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("getBlockNumber", pByte)
}

func (c *Client) GetTransaction(param TransactionParam) (any, error) {
	pByte, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return c.Call("getTransaction", pByte)
}

func (c *Client) Call(method string, param []byte) (any, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	// Generate a snowflake ID.
	id := node.Generate()

	// 创建一个JSON-RPC请求
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method, // 假设的方法名，需要匹配服务器端的方法
		Params:  param,  // 参数，这里是一个JSON数组
		ID:      id.String(),
	}
	// 发送请求并获取响应
	response, err := sendJSONRPCRequest(c.url, request)
	if err != nil {
		return nil, err
	}
	result := response.Result

	switch method {
	case "ord_call":
		return result.pResult()
	case "bestBlockNumber":
		return result.pBestBlockNumber()
		case "getScriptModel":
		return result.pScriptModel()
	case "getTokenModel":
		return result.pTokenModel()
	case "getTokenUri":
		return result.pTokenUri()
	case "getEvents":
		return result.pEvents()
	case "getBlockNumber":
		return result.pBlockNumber()
	case "getTransaction":
		return result.pTransaction()
	}

	return nil, nil
}

func (r JSONResult) pScriptModel() (*Script, error) {
	var sm Script
	if r.Data == nil {
		return nil, errors.New("not find")
	}
	m := r.Data.(map[string]interface{})
	sm.Abi = m["abi"]
	sm.Kip = m["bip"].(string)
	return &sm, nil
}

func (r JSONResult) pTokenModel() (*Token, error) {
	var t Token

	if r.Data == nil {
		return nil, errors.New("token not find")
	}

	m := r.Data.(map[string]interface{})
	t.Name = m["Name"].(string)
	t.Symbol = m["Symbol"].(string)
	t.TotalSupply = fmt.Sprint(m["TotalSupply"])
	if m["Owner"] != nil {
		t.Owner = m["Owner"].(string)
	}
	return &t, nil
}

func (r JSONResult) pTokenUri() (*string, error) {
	if r.Data == nil {
		return nil, errors.New("tokenUri is empty")
	}
	uri := r.Data.(string)
	return &uri, nil
}

func (r JSONResult) pResult() (interface{}, error) {
	if r.Data == nil {
		return nil, nil
	}
	return r.Data, nil
}

func (r JSONResult) pEvents() ([]Event, error) {
	var evenList []Event
	if r.Data == nil {
		return nil, nil
	}
	eInterfaces := r.Data.([]interface{})
	for _, eInter := range eInterfaces {
		m := eInter.(map[string]interface{})
		var event Event
		event.KID = m["kid"].(string)
		event.EHash = m["e_hash"].(string)
		event.TxHash = m["tx_hash"].(string)
		event.Height = int64(m["height"].(float64))
		event.Name = m["name"].(string)

		if m["args"] != nil {
			event.Args = m["args"].(map[string]interface{})
		}

		event.TimeStamp = int64(m["timestamp"].(float64))
		evenList = append(evenList, event)
	}
	return evenList, nil
}

func (r JSONResult) pBestBlockNumber() (int64, error) {
	return int64(r.Data.(float64)), nil
}

func (r JSONResult) pBlockNumber() ([]Transaction, error) {
	var txList []Transaction
	tInterfaces := r.Data.([]interface{})
	for _, tInter := range tInterfaces {
		if tInter == nil {
			continue
		}
		m := tInter.(map[string]interface{})
		var t Transaction
		t.Height = int64(m["height"].(float64))
		t.TxHash = m["tx_hash"].(string)
		t.Sender = m["sender"].(string)
		t.KID = m["kid"].(string)
		t.OP = m["op"].(string)
		t.Input = m["input"].(string)
		t.Out = m["out"]
		t.Logs = m["logs"]
		t.TimeStamp = int64(m["timestamp"].(float64))
		t.Status = int(m["status"].(float64))
		txList = append(txList, t)
	}
	return txList, nil
}

func (r JSONResult) pTransaction() (*Transaction, error) {
	var t Transaction

	if r.Data == nil {
		return nil, errors.New("transaction not find")
	}

	m := r.Data.(map[string]interface{})
	t.Height = int64(m["height"].(float64))
	t.TxHash = m["tx_hash"].(string)
	t.Sender = m["sender"].(string)
	t.KID = m["kid"].(string)
	t.OP = m["op"].(string)
	t.Input = m["input"].(string)
	t.Out = m["out"]
	t.Logs = m["logs"]
	t.TimeStamp = int64(m["timestamp"].(float64))
	t.Status = int(m["status"].(float64))
	return &t, nil
}

func DecodeBytes(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
