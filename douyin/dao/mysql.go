package dao

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"douyin/douyin/common"
)

var Db *gorm.DB

// 配置参数映射结构体

// InitMySql 初始化连接数据库
func InitMySql(config common.Conf) (*sql.DB, error) {
	// 获取yaml配置参数
	conf := config.Mysql
	// 将yaml配置参数拼接成连接数据库的url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.PassWord,
		conf.Url,
		conf.Port,
		conf.DbName,
	)
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db = db
	if err != nil {
		panic(err)
	}
	// 验证数据库连接是否成功，若成功，则无异常
	sqlDb, err := Db.DB()
	if err != nil {
		return nil, err
	}
	return sqlDb, sqlDb.Ping()
}
