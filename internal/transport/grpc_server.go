package transport

import (
	"fmt"
	microservices "git.ramooz.org/ramooz/golang-components/microservice/connection/grpc/middlewares"
	"github.com/mortezakhademan/user-service-sample/internal/config"
	"github.com/mortezakhademan/user-service-sample/internal/repository"
	"github.com/mortezakhademan/user-service-sample/services"
	pbUser "github.com/mortezakhademan/user-service-sample/services/proto/apis-gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

func InitGrpcService(address, port string) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("grpc server ran on %s:%s\n", address, port)
	}

	srv := grpc.NewServer(microservices.Middlewares(
		microservices.GrpcRecovery(),
		microservices.GrpcValidator(),
	))
	userRepo := repository.NewMongoUserRepository(config.DB, "users")
	userService := services.NewUserService(userRepo)
	pbUser.RegisterUserServiceServer(srv, userService)

	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)

	go func() {
		if err := srv.Serve(listener); err != nil {
			config.Logger.Fatal(err)
		}
	}()
	return srv, nil
}
