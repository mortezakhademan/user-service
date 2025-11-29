package transport

import (
	"context"
	"fmt"
	"git.ramooz.org/ramooz/golang-components/microservice/connection/http/handlers"
	microservices "git.ramooz.org/ramooz/golang-components/microservice/connection/http/middlewares"
	"github.com/Ja7ad/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mortezakhademan/user-service-sample/api/swagger"
	"github.com/mortezakhademan/user-service-sample/internal/config"
	pbUser "github.com/mortezakhademan/user-service-sample/services/proto/apis-gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"net/http"
)

const maxMsgSize = 1024 * 1024 * 20

func InitRestService(ctx context.Context, address, port string) (*http.Server, error) {

	grpcClientConn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:%s", config.GetHttpAddress(), config.GetGrpcPort()),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)),
		//grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		//	fmt.Println("Unray inteceptor")
		//	return nil
		//}),
	)
	if err != nil {
		return nil, err
	}

	rMux := runtime.NewServeMux(
		runtime.WithHealthEndpointAt(grpc_health_v1.NewHealthClient(grpcClientConn), "/health"),
		runtime.WithIncomingHeaderMatcher(headers),
		runtime.WithErrorHandler(microservices.ErrorHandler),
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md := make(map[string]string)
			md["http-method"] = r.Method
			return metadata.New(md)
		}),
	)

	if err := registerUserEndpoint(ctx, rMux); err != nil {
		return nil, err
	}

	muxHandlers := http.NewServeMux()

	muxHandlers = handlers.SetRuntimeAsRootHandler(muxHandlers, rMux)
	muxHandlers.Handle("/api-docs/", http.StripPrefix("/api-docs", swaggerui.Handler(swagger.Swagger)))
	srv := &http.Server{
		Handler: microservices.AllowCORS(muxHandlers),
		Addr:    fmt.Sprintf("%s:%s", address, port),
	}

	go func() {
		fmt.Printf("rest server ran on %s:%s\r", address, port)
		if err := srv.ListenAndServe(); err != nil {
			config.Logger.Fatal(err)
		}
	}()

	return srv, nil
}

func registerUserEndpoint(ctx context.Context, mux *runtime.ServeMux) error {
	if err := pbUser.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", config.GetHttpAddress(), config.GetGrpcPort()),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return err
	}
	return nil
}

func headers(key string) (string, bool) {
	switch key {
	default:
		return key, false
	}
}
