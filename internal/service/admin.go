package service

import (
	"context"
	pb "go-leaf/api/admin/v1"
	"go-leaf/internal/logic"
	"go-leaf/internal/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Admin struct {
	pb.UnimplementedAdminServer
	log   *pkg.Helper
	logic *logic.Admin
}

func NewAdmin(logger *pkg.Helper, admin *logic.Admin) *Admin {
	return &Admin{
		logic: admin,
		log:   logger,
	}
}

func (a *Admin) GetHost(ctx context.Context, req *pb.GetHostReq) (*pb.GetHostResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHost not implemented")
}
