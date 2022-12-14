package daemon

import (
	"fmt"
	"os"

	"studyhut/commands"
	"studyhut/controllers"
	"studyhut/models"

	"github.com/astaxie/beego"
	"github.com/kardianos/service"
)

// Daemon 守护进程
type Daemon struct {
	config *service.Config
	errs   chan error
}

// NewDaemon 创建守护进程
func NewDaemon() *Daemon {
	config := &service.Config{
		Name:        "studyhut",                              //服务显示名称
		DisplayName: "studyhut Service",                      //服务名称
		Description: "A document online management program.", //服务描述
		Arguments:   os.Args[1:],
	}
	return &Daemon{
		config: config,
		errs:   make(chan error, 100),
	}
}

// Config 获取守护进程的配置
func (d *Daemon) Config() *service.Config {
	return d.config
}

// Start 开启守护协程
func (d *Daemon) Start(s service.Service) error {
	go func() {
		commands.ResolveCommand(d.config.Arguments)
		commands.RegisterFunction()
		beego.ErrorController(&controllers.ErrorController{})
		models.Init()
		beego.Run()
	}()
	return nil
}

// Stop 停止守护进程
func (d *Daemon) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

// Install 安装守护进程服务
func Install() {
	fmt.Println(os.Args, "---", os.Args[3:])
	d := NewDaemon()
	d.config.Arguments = os.Args[3:]
	s, err := service.New(d, d.config)
	if err != nil {
		beego.Error("Create service error => ", err)
		os.Exit(1)
	}
	err = s.Install()
	if err != nil {
		beego.Error("Install service error:", err)
		os.Exit(1)
	} else {
		beego.Info("Service installed!")
	}
	os.Exit(0)
}

// Uninstall 卸载守护进程服务
func Uninstall() {
	d := NewDaemon()
	s, err := service.New(d, d.config)
	if err != nil {
		beego.Error("Create service error => ", err)
		os.Exit(1)
	}
	err = s.Uninstall()
	if err != nil {
		beego.Error("Install service error:", err)
		os.Exit(1)
	} else {
		beego.Info("Service uninstalled!")
	}
	os.Exit(0)
}

// Restart 重启守护进程服务
func Restart() {
	d := NewDaemon()
	s, err := service.New(d, d.config)
	if err != nil {
		beego.Error("Create service error => ", err)
		os.Exit(1)
	}
	err = s.Restart()
	if err != nil {
		beego.Error("Install service error:", err)
		os.Exit(1)
	} else {
		beego.Info("Service Restart!")
	}
	os.Exit(0)
}
