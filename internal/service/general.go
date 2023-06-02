package service

import (
	"context"
	pb "go-leaf/api/general/v1"
	"go-leaf/internal/logic"
	"go-leaf/internal/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type General struct {
	pb.UnimplementedGeneralServer
	logic *logic.General
	log   *pkg.Helper
}

func NewGeneral(general *logic.General, logger *pkg.Helper) *General {
	return &General{
		logic: general,
		log:   logger,
	}
}

func (g *General) General(ctx context.Context, req *pb.GeneralReq) (*pb.GeneralResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method General not implemented")
}
func (g *General) Parse(ctx context.Context, req *pb.ParseReq) (*pb.ParseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Parse not implemented")
}
