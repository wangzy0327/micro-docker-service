# micro-docker-service
micro-docker-service

## 项目结构介绍
```
项目使用3种平台分别进行构建(amd64、arm64v8、mips64le[loongson])

异构平台轻量级容器虚拟化

项目整体分为四个部分：索引容器、代理容器、状态机、业务容器

索引容器[体积小<10M>、启动快<1s>]负责快速启动发送指令给业务容器，并监听任务执行情况

代理容器接收索引容器的指令分发任务到任务队列中

状态机轮询任务队列负责监听，将任务指定给业务容器

业务容器执行具体的业务逻辑、科学计算，任务完成后回写任务到完成队列中

```

### 索引容器
```
索引容器采用scratch进行镜像构建，代码采用golang语言进行编写，代码目录scratch-client/src/microDockerProject
```
- 静态编译go语言
```
编译为amd64指令集架构,项目名为microDockerProject(在GOPATH路径src目录下),编译后的二进制文件为scratch0
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o scratch0 -ldflags '-s' microDockerProject/

编译为arm64v8指令集架构,项目名为microDockerProject(在GOPATH路径src目录下),编译后的二进制文件为scratch1
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o scratch1 -ldflags '-s' microDockerProject/

编译为mips64le指令集架构(龙芯),项目名为microDockerProject(在GOPATH路径src目录下),编译后的二进制文件为scratch2
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -a -o scratch2 -ldflags '-s' microDockerProject/

```

- 查看go语言支持的平台
```
go tool dist list

```

### 代理容器

```
代理容器采用python2的flask框架编写，代码目录server
不同的指令集架构平台需要替换成不同的源来构建镜像
这里具体介绍一下mips64le架构的基础镜像来源
docker pull huangxg20171010/fedora21-base
docker tag huangxg20171010/fedora21-base fedora21-base

原因分析：
启动docker时，docker进程会创建一个名为docker0的虚拟网桥，用于宿主机与容器之间的通信。
当启动一个docker容器时，docker容器将会附加到虚拟网桥上，容器内的报文通过docker0向外转发。
当docker容器访问宿主机时，如果宿主机服务端口会被防火墙拦截，从而无法连通宿主机，出现No route to host的错误

解决方法：（宿主机上执行）
方法一：关闭防火墙
systemctl stop firewalld
方法二：在防火墙上开放指定端口
firewall-cmd --zone=public --add-port=8081/tcp --permanent
firewall-cmd --reload

```

### 状态机
```
状态机采用python3进行代码编写，脚本为state-machine.py 
状态机监听的队列目录publish
```

### 业务容器
```
业务容器采用java进行代码编写，将代码打成jar包进行业务容器镜像构建；
业务容器进行1-n数字之和计算，n由输入参数指定
```
