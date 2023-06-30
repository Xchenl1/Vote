package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// 连接数据库
func Lianjie() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "CHEN2580", "127.0.0.1:3306", "book_manage")
	var mysqllogger logger.Interface
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{Logger: mysqllogger})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	DB = conn
}
func Gb() {
	GB, _ := DB.DB()
	_ = GB.Close()
}
