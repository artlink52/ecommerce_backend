package user

import (
	"context"
	"fmt"

	pb "github.com/artlink52/ecommerce_backend/proto/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	auth pb.AuthServiceClient
	conn *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to user-service: %w", err)
	}
	return &Client{
		auth: pb.NewAuthServiceClient(conn),
		conn: conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Register(ctx context.Context, email, password string) (int64, error) {
	resp, err := c.auth.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return 0, err
	}
	return resp.GetUserId(), nil
}

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	resp, err := c.auth.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.GetToken(), nil
}
