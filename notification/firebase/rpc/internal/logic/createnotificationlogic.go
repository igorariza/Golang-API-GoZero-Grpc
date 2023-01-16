package logic

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/internal/svc"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/types/notifications/v1alpha1"

	"github.com/zeromicro/go-zero/core/logx"

	firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

type CreateNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotificationLogic {
	return &CreateNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNotificationLogic) CreateNotification(in *v1alpha1.CreateNotificationRequest) (*v1alpha1.CreateNotificationResponse, error) {
	createAt := time.Now()
	fmt.Println("createAt: ", createAt)

	// notification := &v1alpha1.Notification{
	// 	OrganizationId: in.Notification.OrganizationId,
	// 	Title:          in.Notification.Title,
	// 	Description:    in.Notification.Description,
	// 	Type:           in.Notification.Type,
	// 	Status:         in.Notification.Status,
	// 	CreatedAt:      fmt.Sprint(createAt),
	// 	UpdatedAt:      "",
	// 	DeletedAt:      "",
	// }

	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "https://ccp-notification-default-rtdb.firebaseio.com/",
	}

	// fetch service account key
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	fmt.Println("opt: ", opt)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	// create ref at path user_scores/:userId
	ref := client.NewRef("Organizations/" + in.Notification.OrganizationId + "/Notifications/" + in.Notification.Id)

	if err := ref.Set(context.TODO(), map[string]interface{}{
		"OrganizationId": in.Notification.OrganizationId,
		"Title":          in.Notification.Title,
		"Description":    in.Notification.Description,
		"Type":           in.Notification.Type,
		"Status":         in.Notification.Status,
		"CreatedAt":      createAt,
		"UpdatedAt":      "",
		"DeletedAt":      "",
	}); err != nil {
		log.Fatalln("error in setting value: ", err)
	}

	fmt.Println("score added/updated successfully!")

	return &v1alpha1.CreateNotificationResponse{}, nil
}
