package model

import (
	"strconv"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
 dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
*/

func TestConnectMySQL(t *testing.T) {
	dsn := "root:Kirito768168@tcp(127.0.0.1:3307)/im_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(strconv.Itoa(db.CreateBatchSize))
}
