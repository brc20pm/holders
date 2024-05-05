package jsonrpc

type CallParam struct {
	KID    string      `json:"kid"`
	Method string      `json:"method"`
	Params   interface{} `json:"params"`
}

type EventParam struct {
	Number string `json:"number"`
}

type ScriptParam struct {
	KID string `json:"kid"`
}

type TokenParam struct {
	KID string `json:"kid"`
}

type TokenUriParam struct {
	KID     string `json:"kid"`
	TokenId string `json:"tokenId"`
}

type BlockNumberParam struct {
	Number string `json:"number"`
}

type TransactionParam struct {
	Hash string `json:"hash"`
}

type Event struct {
	EHash string `json:"e_hash"` //事件哈希

	Height    int64                  `json:"height"`  //区块高度
	TxHash    string                 `json:"tx_hash"` //交易哈希
	KID       string                 `json:"kid"`     //合约地址
	Name      string                 `json:"name"`    //事件名称
	Args      map[string]interface{} `json:"args"`
	TimeStamp int64                  `json:"timestamp"` //时间戳
}

type Transaction struct {
	TxHash string `json:"tx_hash"` //交易哈希

	Height    int64       `json:"height"`         //区块高度
	Sender    string      `json:"sender"`         //发起人
	KID       string      `json:"kid,omitempty"`  //调用的合约地址
	OP        string      `json:"op"`             //操作识别
	Input     string      `json:"input"`          //交易输入16进制字符串
	Out       interface{} `json:"out,omitempty"`  //交易输出
	Logs      interface{} `json:"logs,omitempty"` //交易产生的日志集合
	TimeStamp int64       `json:"timestamp"`      //时间戳
	Status    int         `json:"status"`         //交易状态 0-失败 1-成功
}

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"totalSupply"`
	Owner       string `json:"owner"`
}

type Script struct {
	Abi interface{} `json:"abi"`
	Kip string      `json:"kip"`
}
