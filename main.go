package main

import (
	"fmt"
	"os"

	"studyhut/commands"
	"studyhut/commands/daemon"
	_ "studyhut/routers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kardianos/service"
)

func main() {
	if len(os.Args) >= 3 && os.Args[1] == "service" {
		if os.Args[2] == "install" {
			daemon.Install()
		} else if os.Args[2] == "remove" {
			daemon.Uninstall()
		} else if os.Args[2] == "restart" {
			daemon.Restart()
		}
	}

	// 注册orm命令行工具
	commands.RegisterCommand()

	// 创建后台程序
	d := daemon.NewDaemon()

	// 创建服务
	s, err := service.New(d, d.Config())
	if err != nil {
		fmt.Println("Create service error => ", err)
		os.Exit(1)
	}

	// 启动服务
	s.Run()
}
