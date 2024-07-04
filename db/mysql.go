package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"holders/conf"
	"holders/models"
	"log"
	"reflect"
	"sync"
)

// 持有表前缀
const HoldTablePrefix = "h_"

// 余额表前缀
const Balance20Prefix = "b2_"
const Balance721Prefix = "b7_"


type MysqlClient struct {
	db    *gorm.DB
	mutex sync.Mutex // 添加互斥锁
}

var MDB *MysqlClient

func init() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:lisp000724@tcp(127.0.0.1:3306)/bits_scanner?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Token{})

	MDB = &MysqlClient{
		db: db,
	}
}

func GetMySQL() *MysqlClient {
	return MDB
}

func CreateTable(model interface{}) (*string, error) {
	var tableName string
	rType := reflect.TypeOf(model)
	switch rType {
	case reflect.TypeOf(&models.Balance20{}):
		balance20 := model.(*models.Balance20)
		if balance20.Kid == "" {
			return nil, errors.New("invalid table name")
		}
		tableName = Balance20Prefix + balance20.Kid
	case reflect.TypeOf(&models.Balance721{}):
		balance721 := model.(*models.Balance721)
		if balance721.Kid == "" {
			return nil, errors.New("invalid table name")
		}
		tableName = Balance721Prefix + balance721.Kid
	case reflect.TypeOf(&models.Wallet{}):
		wallet := model.(*models.Wallet)
		if wallet.Owner == "" {
			return nil, errors.New("invalid table name")
		}
		tableName = HoldTablePrefix + wallet.Owner
	}

	return &tableName, MDB.db.Table(tableName).AutoMigrate(model)
}

func InsertValues(model interface{}) error {
	tableName, err := CreateTable(model)
	if err != nil {
		return err
	}
	return MDB.db.Table(*tableName).Create(model).Error
}

func Token(token models.Token) error {
	return MDB.db.Create(token).Error
}

// 代币转移事务
func Transaction20(transfer20 models.Transfer20) error {
	if transfer20.From == transfer20.To {
		return errors.New("接收地址和发送地址一样")
	}

	if transfer20.Amount < 0 {
		return errors.New("转移数量小于等于0")
	}

	CreateTable(&models.Balance20{
		Kid: transfer20.Kid,
	})

	CreateTable(&models.Wallet{
		Owner: transfer20.To,
	})

	tx := MDB.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//发送地址
	var fBalance models.Balance20
	if transfer20.From != conf.ZeroAddress {
		result := tx.Table(Balance20Prefix+transfer20.Kid).Where("owner = ?", transfer20.From).First(&fBalance)
		//如果owner之前没有数据
		if result.Error == gorm.ErrRecordNotFound {
			return result.Error
		}
		newBalance := fBalance.Amount - transfer20.Amount
		if newBalance > 0 {
			err := tx.Table(Balance20Prefix+transfer20.Kid).Where("owner", transfer20.From).Update("amount", newBalance).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			//删除余额数据
			err := tx.Table(Balance20Prefix+transfer20.Kid).Where("owner", transfer20.From).Delete(nil).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			//删除持有数据
			err = tx.Table(HoldTablePrefix+transfer20.From).Where("kid", transfer20.Kid).Delete(nil).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	var tBalance models.Balance20
	//接收地址
	result := tx.Table(Balance20Prefix+transfer20.Kid).Where("owner = ?", transfer20.To).First(&tBalance)
	//如果owner之前没有数据
	if result.Error == gorm.ErrRecordNotFound {
		//插入余额
		err := tx.Table(Balance20Prefix + transfer20.Kid).Create(&models.Balance20{
			Amount: transfer20.Amount,
			Owner:  transfer20.To,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		//插入持有
		err = tx.Table(HoldTablePrefix + transfer20.To).Create(&models.Wallet{
			Kid: transfer20.Kid,
			Bip: 20,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		//更新持有
		newBalance := tBalance.Amount + transfer20.Amount
		err := tx.Table(Balance20Prefix+transfer20.Kid).Where("owner", transfer20.To).Update("amount", newBalance).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// NFT转移事务
func Transaction721(transfer721 models.Transfer721) error {
	if transfer721.From == transfer721.To {
		return errors.New("接收地址和发送地址一样")
	}

	CreateTable(&models.Balance721{
		Kid: transfer721.Kid,
	})

	CreateTable(&models.Wallet{
		Owner: transfer721.To,
	})

	tx := MDB.db.Begin()

	//先查询有没有
	var count int64
	err := tx.Table(Balance721Prefix+transfer721.Kid).Where("token_id = ?", transfer721.TokenId).Count(&count).Error
	if count == 0 {
		//保存tokenId所有者
		err := tx.Table(Balance721Prefix + transfer721.Kid).Create(&models.Balance721{
			Owner:   transfer721.To,
			TokenId: transfer721.TokenId,
			Data:    transfer721.Data,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		//更新tokenId所有者
		err := tx.Table(Balance721Prefix+transfer721.Kid).Where("token_id", transfer721.TokenId).Update("owner", transfer721.To).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	var tHold int64
	if transfer721.From != conf.ZeroAddress {
		err = tx.Table(Balance721Prefix+transfer721.Kid).Where("owner", transfer721.From).Count(&tHold).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if tHold == 0 {
			//删除持有数据
			err = tx.Table(HoldTablePrefix+transfer721.From).Where("kid", transfer721.Kid).Delete(nil).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	//接收地址
	//查询更新有持有数量
	err = tx.Table(Balance721Prefix+transfer721.Kid).Where("owner", transfer721.To).Count(&tHold).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//如果更新后为1,则第一次持有
	if tHold == 1 {
		//需要记录持有
		err = tx.Table(HoldTablePrefix + transfer721.To).Create(&models.Wallet{
			Kid: transfer721.Kid,
			Bip: 721,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// 钱包持有数据
func FindWalletHold(owner string) (map[string]interface{}, error) {
	var tokens []models.Wallet
	err := MDB.db.Table(HoldTablePrefix + owner).Find(&tokens).Error
	if err != nil {
		return nil, err
	}

	var hMap = make(map[string]interface{})

	var hold20s []models.Hold
	var hold721s []models.Hold

	for _, token := range tokens {
		switch token.Bip {
		case 20:
			var h20 models.Hold
			err := MDB.db.Table(Balance20Prefix+token.Kid).Where("owner", owner).Find(&h20).Error
			if err != nil {
				log.Println(err)
				continue
			}
			err = MDB.db.Table("tokens").Where("kid", token.Kid).Find(&h20).Error
			if err != nil {
				log.Println(err)
				continue
			}

			hold20s = append(hold20s, h20)
		case 721:
			var h721 models.Hold

			var count int64
			err := MDB.db.Table(Balance721Prefix+token.Kid).Where("owner", owner).Count(&count).Error
			if err != nil {
				log.Println(err)
				continue
			}

			h721.Amount = fmt.Sprint(count)

			err = MDB.db.Table("tokens").Where("kid", token.Kid).Find(&h721).Error
			if err != nil {
				log.Println(err)
				continue
			}
			hold721s = append(hold721s, h721)
		}
	}

	hMap["t20"] = hold20s
	hMap["t721"] = hold721s

	return hMap, nil
}

// 持有的tokenId列表
func FindTokenIds(kid, owner string) (tokenIds []models.TokenIds, err error) {
	err = MDB.db.Table(Balance721Prefix+kid).Select("token_id,data").Where("owner", owner).Find(&tokenIds).Error
	if err != nil {
		return nil, err
	}
	return tokenIds, nil
}

// 获取持有分布
func FindDist(kid string, is20 bool) ([]models.Dist, error) {
	var (
		tableName string
		query     string
	)

	var err error

	var distList []models.Dist

	//如果是代币
	if is20 {
		tableName = Balance20Prefix + kid
		err = MDB.db.Table(tableName).Order("amount desc").Limit(100).Find(&distList).Error
	} else {
		tableName = Balance721Prefix + kid
		query = "count(`owner`) as amount,owner"
		err = MDB.db.Table(tableName).Select(query).Group("owner").Order("amount desc").Limit(100).Find(&distList).Error
	}
	if err != nil {
		return nil, err
	}
	return distList, nil
}

// 查询代币
func FindToken(kid string) (models.Token, error) {
	var token models.Token
	err := MDB.db.Table("tokens").Where("kid", kid).Find(&token).Error
	if err != nil {
		return models.Token{}, err
	}
	return token, nil
}
