// Code generated by goctl. DO NOT EDIT.
// Source: notification_api.proto

package notificationapiservice

import (
	"context"

	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/types/notifications/v1alpha1"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateNotificationRequest  = v1alpha1.CreateNotificationRequest
	CreateNotificationResponse = v1alpha1.CreateNotificationResponse
	Notification               = v1alpha1.Notification

	NotificationAPIService interface {
		CreateNotification(ctx context.Context, in *CreateNotificationRequest, opts ...grpc.CallOption) (*CreateNotificationResponse, error)
	}

	defaultNotificationAPIService struct {
		cli zrpc.Client
	}
)

func NewNotificationAPIService(cli zrpc.Client) NotificationAPIService {
	return &defaultNotificationAPIService{
		cli: cli,
	}
}

func (m *defaultNotificationAPIService) CreateNotification(ctx context.Context, in *CreateNotificationRequest, opts ...grpc.CallOption) (*CreateNotificationResponse, error) {
	client := v1alpha1.NewNotificationAPIServiceClient(m.cli.Conn())
	return client.CreateNotification(ctx, in, opts...)
}
