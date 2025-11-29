package connection

import (
	"context"
	"git.ramooz.org/ramooz/golang-components/microservice/helper"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// GracefullyShutdownServer shutdown grpc and http server on interrupt signal
func GracefullyShutdownServer(ctx context.Context, grpcServer *grpc.Server, httpServer *http.Server) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		helper.Log.Warnf("os signal %s received", s.String())
	}

	grpcServer.GracefulStop()
	return httpServer.Shutdown(ctx)
}
