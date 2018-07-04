package grpc

import (
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"golang.org/x/net/context"
	pb "github.com/hanakoa/alpaca/services/mfa/pb"
	"github.com/hanakoa/alpaca/services/mfa/models"
	"database/sql"
	"github.com/google/uuid"
	"time"
	"math/rand"
	"github.com/hanakoa/alpaca/services/mfa/services"
)

type GrpcServer struct {
	DB            *sql.DB
	Port          int
	TwilioService services.TwilioSvc
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	start := time.Now()
	s1 := rand.NewSource(time.Now().UnixNano())
	log.Println("Seed end time:", time.Since(start))
	r1 := rand.New(s1)
	for i := range b {
		b[i] = letters[r1.Intn(len(letters))]
	}
	return string(b)
}

func (service *GrpcServer) Send2FACode(ctx context.Context, in *pb.Send2FACodeRequest) (*pb.Send2FACodeResponse, error) {
	var mfaCode *models.MFACode
	if id, err := uuid.Parse(in.ResetCode); err != nil {
		return nil, err
	} else {
		mfaCode = &models.MFACode{
			ID:         id,
			Code:       randSeq(6),
			Created:    time.Now(),
			Expiration: time.Now().Add(time.Minute * 30),
			Usable:     true,
			Used:       false,
			AccountID:   in.AccountId}
	}

	if tx, err := service.DB.Begin(); err != nil {
		return nil, err
	} else {
		mfaCode.Create(tx)
		// TODO in rabbit mode, we should enqueue the code before persisting
		// in non-rabbit mode, we just fail if the sms fails
		// TODO send phone number with request
		service.TwilioService.SendSms("555-555-5555", "hello from golang")
		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}

	return &pb.Send2FACodeResponse{ResetCode: in.ResetCode, AccountId: in.AccountId}, nil
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
	pb.RegisterSend2FACodeServiceServer(server, service)

	// Register reflection service on gRPC server.
	reflection.Register(server)
	log.Println("Registered grpc services...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
