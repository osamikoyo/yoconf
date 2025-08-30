package grpcserver

import (
	"context"

	"github.com/osamikoyo/yoconf/core"
	"github.com/osamikoyo/yoconf/models"
	"github.com/osamikoyo/yoconf/pb"
)

type GRPCServer struct {
	pb.UnimplementedYoConfServer
	core *core.Core
}

func NewGRPCServer(core *core.Core) *GRPCServer {
	return &GRPCServer{
		core: core,
	}
}

func (s *GRPCServer) CreateChunk(ctx context.Context, chunk *pb.Chunk) (*pb.Resp, error) {
	if err := s.core.NewConfig(&models.Chunk{
		Project: chunk.Project,
		Data:    chunk.Data,
		Version: int(chunk.Version),
		InUse:   chunk.InUse,
	}); err != nil {
		return &pb.Resp{
			Message: err.Error(),
		}, err
	}

	return &pb.Resp{
		Message: "ok",
	}, nil
}

func (s *GRPCServer) RollOn(ctx context.Context, req *pb.RollOnRequest) (*pb.Resp, error) {
	if err := s.core.RollOn(req.Project, int(req.Version)); err != nil {
		return &pb.Resp{
			Message: err.Error(),
		}, err
	}

	return &pb.Resp{
		Message: "ok",
	}, nil
}

func (s *GRPCServer) DeleteChunk(ctx context.Context, req *pb.DeleteRequest) (*pb.Resp, error) {
if err := s.core.DeleteChunk(req.Project, int(req.Version)); err != nil {
		return &pb.Resp{
			Message: err.Error(),
		}, err
	}

	return &pb.Resp{
		Message: "ok",
	}, nil
}
