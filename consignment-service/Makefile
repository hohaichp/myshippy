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
    # go install
	GOOS=linux GOARCH=amd64 go build
    # 根据当前目录下的 Dockerfile 生成名为 consignment-service 的镜像  代码中host.docker.internal替换127.0.0.1
	docker build -t consignment-service .
run:
    # 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
    # 可添加 -d 参数将微服务放到后台运行
	#docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns consignment-service
    docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=consul consignment-service
    #docker run -p 50051:50051 consignment-service