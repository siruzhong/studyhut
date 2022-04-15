# studyhut

## 站点介绍
学习小屋，一个基于Beego开发的在线IT技术资源整合、在线学习、交流分享的站点。每一名用户都是内容的创造者，分享你认为优质的资源，让我们一起学习！一起进步！

**站点地址**：http://studyhut.cn/（目前没有升级为`https`协议）

![image-20220415132404843](https://bareth-1305674339.cos.ap-hongkong.myqcloud.com/img/image-20220415132404843.png)

## 安装教程

1、下载源代码

```shell
git clone https://gitee.com/bareth/studyhut.git
```

2、修改数据库配置 `conf/app.conf`

```shell
vim conf/app.conf
```

3、安装数据库

```shell
./main install
```

4、启动项目

```shell
go run main.go
```

然后访问本机的80端口即可访问



