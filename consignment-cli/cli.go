package main

import (
	pb "github.com/hohaichp/myshippy/consignment-cli/proto/consignment"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, err
}

func main() {

	// Register consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		// options.Addrs = []string {"127.0.0.1:8500"}
		options.Addrs = []string {"host.docker.internal:8500"}
	})

	// Create service
	server := micro.NewService(
		// 必须和 consignment.proto 中的 package 一致
		micro.Name("shippy.cli.consignment"),
		micro.Version("latest"),
		micro.Registry(reg),
	)
	// 解析命令行参数
	server.Init()

	client := pb.NewShippingService("shippy.service.consignment", server.Client())

	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
