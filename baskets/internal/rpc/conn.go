package rpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	defer func() {
		go func() {
			<-ctx.Done()
			if err := conn.Close(); err != nil {
				slog.Warn("connection close unexpectedly")
			}
		}()
	}()

	return conn, nil
}
