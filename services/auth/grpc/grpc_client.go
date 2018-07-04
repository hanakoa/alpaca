package grpc

import (
	"golang.org/x/net/context"
	pb "github.com/hanakoa/alpaca/services/auth/pb"
	"fmt"
	grpcUtils "github.com/kevinmichaelchen/my-go-utils/grpc"
	"time"
)

type AuthClient = pb.AccountServiceClient
type PassClient = pb.ResetPasswordServiceClient

func NewPassClient(host string, port int) PassClient {
	conn := grpcUtils.InitGrpcConn(fmt.Sprintf("%s:%d", host, port), 3, time.Second*5)
	return pb.NewResetPasswordServiceClient(conn)
}

func GetAccountIDForEmailAddress(client AuthClient, emailAddress string) (int64, error) {
	request := &pb.GetAccountRequest{EmailAddress: emailAddress}
	if response, err := client.GetAccount(context.Background(), request); err != nil {
		return 0, err
	} else {
		return response.AccountId, nil
	}
}

func ResetPassword(client PassClient, accountID int64, newPassword string) error {
	// We allow nil clients so that unit tests can pass nil to effectively disable gRPC.
	if client != nil {
		request := &pb.ResetPasswordRequest{AccountId: accountID, NewPassword: newPassword}
		_, err := client.ResetPassword(context.Background(), request)
		return err
	}
	return nil
}