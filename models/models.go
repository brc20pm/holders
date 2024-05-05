package models


type Transfer20 struct {
	Kid    string  `json:"kid"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type Transfer721 struct {
	Kid     string      `json:"kid"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	TokenId interface{} `json:"tokenId"`
	Data    string      `json:"data"`
}

type Balance20 struct {
	Kid    string  `gorm:"-"`
	Amount float64 `json:"amount gorm:"type:float64"`
	Owner  string  `json:"owner" gorm:"uniqueIndex"`
}

type Balance721 struct {
	Kid     string      `gorm:"-"`
	TokenId interface{} `json:"tokenId" gorm:"uniqueIndex;type:string"`
	Owner   string      `json:"-"`
	Data    string      `json:"data"`
}

type TokenIds struct {
	TokenId string `json:"tokenId"`
	Data    string `json:"data"`
}

type Wallet struct {
	Owner string `gorm:"-"`
	Kid   string `json:"kid" gorm:"uniqueIndex"`
	Bip   int    `json:"bip"`
}

type Token struct {
	Kid         string `json:"kid" gorm:"uniqueIndex"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"totalSupply"`
	Other       string `json:"other"`
}

type Hold struct {
	Kid    string `json:"kid"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

type Dist struct {
	Owner  string `json:"owner"`
	Amount string `json:"amount"`
}

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  interface{} `json:"msg"`
}
