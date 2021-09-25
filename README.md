# myshippy 
myshippy

#makefile

##consignment-cli - makefile
```
build:
    # 告知 Go 编译器生成二进制文件的目标环境：amd64 CPU 的 Linux 系统, 以下是windows版本
    # SET CGO_ENABLED=0
    # SET GOOS=linux
    # SET GOARCH=amd64
    # go build
	GOOS=linux GOARCH=amd64 go build
    # 根据当前目录下的 Dockerfile 生成名为 consignment-service 的镜像  代码中host.docker.internal替换127.0.0.1
	docker build -t consignment-cli .
run:
    # 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
    # 可添加 -d 参数将微服务放到后台运行
    #docker run -e MICRO_REGISTRY=mdns consignment-cli
    docker run -e MICRO_REGISTRY=consul consignment-cli
```
---
##consignment-service - makefile
```
build:
	# 不再使用 grpc 插件
	protoc -I. --go_out=plugins=micro:D:/workspaces/goproject/shippy/consignment-service proto/consignment/consignment.proto
    # v3版本
	#protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. greeter.proto
	protoc --proto_path=./proto/consignment/ --micro_out=D:/workspaces/goproject/shippy/consignment-service --go_out=D:/workspaces/goproject/shippy/consignment-service consignment.proto

    # 告知 Go 编译器生成二进制文件的目标环境：amd64 CPU 的 Linux 系统, 以下是windows版本
    # SET CGO_ENABLED=0
    # SET GOOS=linux
    # SET GOARCH=amd64
    # go build
	GOOS=linux GOARCH=amd64 go build
    # 根据当前目录下的 Dockerfile 生成名为 consignment-service 的镜像  代码中host.docker.internal替换127.0.0.1
	docker build -t consignment-service .
run:
    # 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
    # 可添加 -d 参数将微服务放到后台运行
	#docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns consignment-service
    docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=consul consignment-service
    #docker run -p 50051:50051 consignment-service
```
---
##vessel-service - makefile
```
build:
	# 不再使用 grpc 插件 v2版本
	protoc -I. --go_out=plugins=micro:D:/workspaces/goproject/shippy/vessel-service proto/vessel/vessel.proto
    # v3版本
	#protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. greeter.proto
	protoc --proto_path=./proto/vessel/ --micro_out=D:/workspaces/goproject/shippy/vessel-service --go_out=D:/workspaces/goproject/shippy/vessel-service vessel.proto

    # 告知 Go 编译器生成二进制文件的目标环境：amd64 CPU 的 Linux 系统, 以下是windows版本
    # SET CGO_ENABLED=0
    # SET GOOS=linux
    # SET GOARCH=amd64
    # go build
    # go install
	GOOS=linux GOARCH=amd64 go build
    # 根据当前目录下的 Dockerfile 生成名为 consignment-service 的镜像  代码中host.docker.internal替换127.0.0.1
	docker build -t vessel-service .
run:
    # 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
    # 可添加 -d 参数将微服务放到后台运行
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns vessel-service
run:
    # 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
    # 可添加 -d 参数将微服务放到后台运行
	#docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns vessel-service
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=consul vessel-service
	#docker run -p 50052:50051 vessel-service
```
---
# consul

获取consul镜像:

`docker pull consul`

启动consul命令:

`docker run -d --name=dev-consul -p 8500:8500 consul`
    
`docker run -d -e CONSUL_BIND_INTERFACE=eth0 consul agent -dev -join=127.0.0.1`

`docker exec -t dev-consul consul members`

--go micro - plugin - consul
```
go get -u github.com/asim/go-micro/plugins/registry/consul/v3
go get -u github.com/asim/go-micro/v3/registry
```





