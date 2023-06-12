package grpc_client

import (
	"context"

	"github.com/majidbl/discount/internal/interceptors"
	walletService "github.com/majidbl/discount/proto/wallet"
)

func NewWalletGrpcClient(
	ctx context.Context,
	port string, im interceptors.InterceptorManager,
) (walletService.WalletServiceClient, func() error, error) {
	grpcServiceConn, err := NewGrpcServiceConn(ctx, port, im)
	if err != nil {
		return nil, nil, err
	}

	serviceClient := walletService.NewWalletServiceClient(grpcServiceConn)

	return serviceClient, grpcServiceConn.Close, nil
}
