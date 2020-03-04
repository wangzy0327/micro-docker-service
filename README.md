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
编译为amd64指令集架构,项目名为scratchHadoopProject(在GOPATH路径src目录下),编译后的二进制文件为scratchHadoop
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o scratchHadoop -ldflags '-s' scratchHadoopProject/

编译为mips64le指令集架构(龙芯),项目名为scratchHadoopProject(在GOPATH路径src目录下),编译后的二进制文件为scratchHadoop0
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -a -o scratchHadoop0 -ldflags '-s' scratchHadoopProject/

编译为arm64指令集架构,项目名为scratchHadoopProject(在GOPATH路径src目录下),编译后的二进制文件为scratchHadoop1
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o scratchHadoop1 -ldflags '-s' scratchHadoopProject/

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