package main

import (
	pb "github.com/hohaichp/myshippy/vessel-service/proto/vessel"
	"context"
	"log"
	"github.com/pkg/errors"
	"github.com/micro/go-micro/v2"
)

var count int

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

// 接口实现
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	// 选择最近一条容量、载重都符合的货轮
	for _, v := range repo.vessels {
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("No vessel can't be use")
}

// 定义货船服务
type service struct {
	repo Repository
}

// 实现服务端
func (s *service) FindAvailable(ctx context.Context, spec *pb.Specification, resp *pb.VesselResponse) error {
	// 调用内部方法查找
	v, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	count++
	log.Printf("FindAvailable succeed, count:%d, Id: %s, Capacity: %d, MaxWeight: %d\n",
		count, v.Id, v.Capacity, v.MaxWeight)
	resp.Vessel = v
	return nil
}

func main() {
	// 停留在港口的货船，先写死
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}
	server := micro.NewService(
		micro.Name("shippy.service.vessel"),
		micro.Version("latest"),
	)
	server.Init()

	// 将实现服务端的 API 注册到服务端
	pb.RegisterVesselServiceHandler(server.Server(), &service{repo})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
