package dao

import (
	"database/sql"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// 配置参数映射结构体
type Conf struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"username"`
	PassWord string `yaml:"password"`
	DbName   string `yaml:"database"`
	Port     string `yaml:"port"`
}

// 获取配置参数数据
func (c *Conf) getConf() *Conf {
	// 读取resources/application.yaml文件
	yamlFile, err := os.ReadFile("conf/application.yaml")
	// 若出现错误，打印错误提示
	if err != nil {
		fmt.Println(err.Error())
	}
	// 将读取的字符串转换成结构体conf
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

// InitMySql 初始化连接数据库
func InitMySql() (*sql.DB, error) {
	var c Conf
	// 获取yaml配置参数
	conf := c.getConf()
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
