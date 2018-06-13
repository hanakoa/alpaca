package grpc

import (
	"golang.org/x/net/context"
	pb "github.com/hanakoa/alpaca/mfa/pb"
	"fmt"
	grpcUtils "github.com/kevinmichaelchen/my-go-utils/grpc"
	"time"
	"github.com/google/uuid"
)

type MFAClient = pb.Send2FACodeServiceClient

func NewMFAClient(host string, port int) MFAClient {
	conn := grpcUtils.InitGrpcConn(fmt.Sprintf("%s:%d", host, port), 3, time.Second*5)
	return pb.NewSend2FACodeServiceClient(conn)
}

func Send2FACode(client MFAClient, personID int64, resetCode uuid.UUID) error {
	request := &pb.Send2FACodeRequest{ResetCode: resetCode.String(), PersonId: personID}
	_, err := client.Send2FACode(context.Background(), request)
	return err
}