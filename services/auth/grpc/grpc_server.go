package grpc

import (
	"fmt"
	"log"
	pb "github.com/hanakoa/alpaca/services/auth/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"database/sql"
	"github.com/hanakoa/alpaca/services/auth/models"
	"github.com/hanakoa/alpaca/services/auth/services"
	"gopkg.in/guregu/null.v3"
)

type GrpcServer struct {
	DB              *sql.DB
	Port            int
	PasswordService services.PasswordService
}

func (service *GrpcServer) ResetPassword(ctx context.Context, in *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if tx, err := service.DB.Begin(); err != nil {
		return nil, err
	} else {
		p := &models.Password{PasswordText: null.StringFrom(in.NewPassword), AccountID: in.AccountId}
		if _, err := service.PasswordService.UpdatePasswordHelper(tx, p, in.AccountId); err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &pb.ResetPasswordResponse{AccountId: in.AccountId}, nil
	}
}

func (service *GrpcServer) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	log.Printf("Looking up account for email address: %s", in.EmailAddress)
	e := &models.EmailAddress{EmailAddress: in.EmailAddress}
	if err := e.GetEmailAddressByEmailAddress(service.DB); err != nil {
		return &pb.GetAccountResponse{AccountId: 0}, err
	} else {
		log.Printf("Found accountId: %d", e.AccountID)
		return &pb.GetAccountResponse{AccountId: e.AccountID}, nil
	}
}

func (service *GrpcServer) Run() {
	address := fmt.Sprintf(":%d", service.Port)
	log.Printf("Listening on %s\n", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Starting grpc server...")
	server := grpc.NewServer()

	// Register our services
	pb.RegisterAccountServiceServer(server, service)
	pb.RegisterResetPasswordServiceServer(server, service)

	// Register reflection service on gRPC server.
	reflection.Register(server)
	log.Println("Registered grpc services...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
