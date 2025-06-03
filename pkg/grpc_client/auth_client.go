package grpc_client

import (
	"context"
	"fmt"
	"log"
	"time"

	authv1 "github.com/SeiFlow-3P2/auth_service/pkg/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient interface {
	GetTelegramIDByUserID(ctx context.Context, userID string) (string, error)
	Close() error
}

type authServiceClient struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func NewAuthServiceClient(ctx context.Context, targetAddress string) (AuthServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		targetAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %w", targetAddress, err)
	}

	log.Printf("Successfully connected to auth service at %s", targetAddress)
	client := authv1.NewAuthServiceClient(conn)
	return &authServiceClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *authServiceClient) GetTelegramIDByUserID(ctx context.Context, userID string) (string, error) {
	req := &authv1.GetUserInfoRequest{UserId: userID}

	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetUserInfo(callCtx, req)
	if err != nil {
		return "", fmt.Errorf("auth service GetUserInfo failed: %w", err)
	}

	if resp.GetUserInfo() == nil {
		return "", fmt.Errorf("auth service returned nil UserInfo for user_id: %s", userID)
	}

	telegramID := resp.GetUserInfo().GetTelegramId()
	if telegramID == "" {
		log.Printf("Warning: Telegram ID not found for user_id: %s", userID)

	}

	return telegramID, nil
}

func (c *authServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
