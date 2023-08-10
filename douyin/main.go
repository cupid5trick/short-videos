package main

import (
	"time"

	"douyin/douyin/dao"
	"douyin/douyin/model"
	"douyin/douyin/routes"
)

func main() {
	// 连接数据库
	sqlDB, err := dao.InitMySql()
	if err != nil {
		panic(err)
	}
	// 程序退出关闭数据库连接

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 绑定模型
	// r := gin.Default()
	dao.Db.AutoMigrate(&model.User{})
	dao.Db.AutoMigrate(&model.Video{})
	dao.Db.AutoMigrate(&model.Comment{})
	dao.Db.AutoMigrate(&model.Favorite{})
	dao.Db.AutoMigrate(&model.Following{})
	dao.Db.AutoMigrate(&model.Followers{})
	// 注册路由
	r := routes.InitRouter()
	// 启动端口为8080的项目
	errRun := r.Run(":8000")
	if errRun != nil {
		return
	}
}
