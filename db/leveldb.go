package db

import (
	"log"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelClient struct {
	DB    *leveldb.DB
	mutex sync.Mutex // 添加互斥锁
}

var LDB *LevelClient

func init() {
	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
		log.Fatal(err)
	}

	LDB = &LevelClient{DB: db}
}

// 打开或创建LevelDB数据库
func GetLevelDB() *LevelClient {
	return LDB
}

// 批量添加
func (c *LevelClient) Batch(batch *leveldb.Batch) error {
	c.mutex.Lock()         // 在操作前加锁
	defer c.mutex.Unlock() // 确保在操作后解锁
	err := c.DB.Write(batch, nil)
	return err
}

// 存储数据
func (c *LevelClient) Put(key, value []byte) error {
	c.mutex.Lock()         // 在操作前加锁
	defer c.mutex.Unlock() // 确保在操作后解锁

	err := c.DB.Put(key, value, nil)
	return err
}

// 获取数据
func (c *LevelClient) Get(key string) ([]byte, error) {
	c.mutex.Lock()         // 在操作前加锁
	defer c.mutex.Unlock() // 确保在操作后解锁

	data, err := c.DB.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 删除数据
func (c *LevelClient) Delete(key string) error {
	c.mutex.Lock()         // 在操作前加锁
	defer c.mutex.Unlock() // 确保在操作后解锁

	err := c.DB.Delete([]byte(key), nil)
	return err
}

// key是否存在
func (c *LevelClient) Has(key string) bool {
	c.mutex.Lock()         // 在操作前加锁
	defer c.mutex.Unlock() // 确保在操作后解锁
	return c.Has(key)
}

// 获取最新区块
func FistNumber(chainId string) uint64 {
	n, err := LDB.Get(chainId)
	if err != nil {
		return 0
	}
	number, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		return 0
	}
	return number
}

// 写入最新区块
func WriteNumber(chainId string, number string) error {
	err := LDB.Put([]byte(chainId), []byte(number))
	return err
}

func PutTokenExits(kid string) error {
	return LDB.Put([]byte(kid), []byte("bool"))
}

func GetTokenExits(kid string) bool {
	_, err := LDB.Get(kid)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return false
		}
	}
	return true
}

func PutTokenUriExits(kid, tokenId string) error {
	return LDB.Put([]byte(kid+tokenId), []byte("bool"))
}

func GetTokenUriExits(kid, tokenId string) bool {
	_, err := LDB.Get(kid + tokenId)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return false
		}
	}
	return true
}
