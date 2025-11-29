package serviceInfo

import (
	service "git.ramooz.org/ramooz/pb/apis-gen/services/user/v2"
	"google.golang.org/grpc"
)

var serviceInfoClient service.ServiceInfoServiceClient
var userServiceConnection *grpc.ClientConn

func SetServiceInfoClientConnection(conn *grpc.ClientConn) {
	userServiceConnection = conn

	serviceInfoClient = service.NewServiceInfoServiceClient(userServiceConnection)
	//userClient = NewUserServiceClient(userServiceConnection)
}

func GetServiceInfoClient() service.ServiceInfoServiceClient {
	return serviceInfoClient
}
