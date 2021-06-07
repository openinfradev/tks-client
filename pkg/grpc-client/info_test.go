package grpcclient_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	gc "github.com/sktelecom/tks-client/pkg/grpc-client"
	pb "github.com/sktelecom/tks-proto/pbgo"
	mock "github.com/sktelecom/tks-proto/pbgo/mock"
)

func TestNewClient(t *testing.T) {
	clientConn, serviceClient, err := gc.CreateClientsObject("localhost", 9110, false, "")
	if err != nil {
		t.Errorf("Unexpected error while creating new client %s", err)
	}
	infoClient := gc.New(clientConn, serviceClient)
	defer infoClient.Close()
}

func TestCreateCSPInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	mockClient := mock.NewMockInfoServiceClient(ctrl)
	infoClient := gc.New(nil, mockClient)

	defer func() {
		cancel()
	}()
	mockClient.EXPECT().CreateCSPInfo(
		gomock.Any(),
		&pb.CreateCSPInfoRequest{
			ContractId: "xxxx-xxxxxxx-xxxx-xxxx",
			CspName:    "aws",
			Auth:       "auth",
		},
	).Return(&pb.IDResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Id:    "a254a66e-7225-4527-bf4c-9b5494c99b37",
	}, nil)
	res, err := infoClient.CreateCSPInfo(ctx, "xxxx-xxxxxxx-xxxx-xxxx", "aws", "auth")
	if err != nil {
		t.Errorf("an unexpected error while creating csp info %s", err)
	}
	if res.Code != pb.Code_OK_UNSPECIFIED {
		t.Error("Not expected response code:", res.Code)
	}
	if res.GetId() == "" {
		t.Error("CspId is empty.")
	}
}
