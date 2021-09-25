package main

import (
	pb "github.com/hohaichp/myshippy/consignment-service/proto/consignment"
	vesselPb "github.com/hohaichp/myshippy/consignment-service/proto/vessel"
	"context"
	"log"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
)

var count int

//
// 仓库接口
//
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
	GetAll() []*pb.Consignment                                   // 获取仓库中所有的货物
}

//
// 我们存放多批货物的仓库，实现了 IRepository 接口
//
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

//
// 定义微服务
//
type service struct {
	repo Repository
	// consignment-service 作为客户端调用 vessel-service 的函数
	vesselClient vesselPb.VesselService
}

//
// 实现 consignment.pb.go 中的 ShippingServiceHandler 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
// func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {

	// 检查是否有适合的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}
	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	count++
	log.Printf("CreateConsignment succeed, count:%d\n", count)
	resp.Created = true
	resp.Consignment = consignment
	return nil
}

// 获取目前所有托运的货物
// func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	resp.Consignments = allConsignments
	return nil
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
		micro.Name("shippy.service.consignment"),
		micro.Version("latest"),
		micro.Registry(reg),
		)
	// 解析命令行参数
    server.Init()
	repo := Repository{}

	// 作为 vessel-service 的客户端
	vClient := vesselPb.NewVesselService("shippy.service.vessel", server.Client())

	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}