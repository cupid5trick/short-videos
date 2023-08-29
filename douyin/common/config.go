package common

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ConfPath = "conf/application.yaml"
)

var (
	APPConfig Conf
)

type Conf struct {
	Mysql Mysql
	Fs    Fs
}

type Mysql struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"username"`
	PassWord string `yaml:"password"`
	DbName   string `yaml:"database"`
	Port     string `yaml:"port"`
}

type Fs struct {
	Type string
	URI  string `yaml:"uri"`
}

// LoadConfig 获取配置参数数据
func LoadConfig() *Conf {
	// 读取resources/application.yaml文件
	yamlFile, err := os.ReadFile(ConfPath)
	// 若出现错误，打印错误提示
	if err != nil {
		fmt.Println(err.Error())
	}

	// 将读取的字符串转换成结构体conf
	err = yaml.Unmarshal(yamlFile, &APPConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &APPConfig
}
