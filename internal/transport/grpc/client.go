package grpc

import (
	"context"
	"fmt"

	logger "github.com/radmirid/grpc-logger/pkg/domain"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn         *grpc.ClientConn
	loggerClient logger.GRPCServiceClient
}

func NewClient(port int) (*Client, error) {
	var conn *grpc.ClientConn

	addr := fmt.Sprintf(":%d", port)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:         conn,
		loggerClient: logger.NewGRPCServiceClient(conn),
	}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) SendLogRequest(ctx context.Context, req logger.LogItem) error {
	action, err := logger.ToPbAction(req.Action)
	if err != nil {
		return err
	}

	entity, err := logger.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}

	_, err = c.loggerClient.Log(ctx, &logger.LogRequest{
		Action:    action,
		Entity:    entity,
		EntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})

	return err
}
