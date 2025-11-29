package helper

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func BuildTarget(ipAddress, port string) string {
	if port != "" {
		return fmt.Sprintf("%s:%s", ipAddress, port)
	} else {
		return fmt.Sprintf("%s", ipAddress)
	}
}

// GrpcContextHeader create GRPC context header
func GrpcContextHeader(header map[string]string) context.Context {
	md := metadata.New(header)
	return metadata.NewOutgoingContext(context.Background(), md)
}

// GrpcIncomingContextHeader create GRPC incoming context header
func GrpcIncomingContextHeader(header map[string]string) context.Context {
	md := metadata.New(header)
	return metadata.NewIncomingContext(context.Background(), md)
}

// GrpcMultiContextHeader append multi GRPC context in a context container
func GrpcMultiContextHeader(appendCtx ...context.Context) context.Context {
	multiMD := make([]metadata.MD, 0)
	for i := range appendCtx {
		if ctx, ok := metadata.FromOutgoingContext(appendCtx[i]); ok {
			multiMD = append(multiMD, ctx)
		}
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Join(multiMD...))
}

// GrpcMultiContextIncomingHeader append multi GRPC incoming context in a context container
func GrpcMultiContextIncomingHeader(appendCtx ...context.Context) context.Context {
	multiMD := make([]metadata.MD, 0)
	for i := range appendCtx {
		if ctx, ok := metadata.FromIncomingContext(appendCtx[i]); ok {
			multiMD = append(multiMD, ctx)
		}
	}
	return metadata.NewIncomingContext(context.Background(), metadata.Join(multiMD...))
}

// GrpcContextHeaderExtractor extract header from context
func GrpcContextHeaderExtractor(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}
	if len(md.Get(header)) != 0 {
		return md.Get(header)[0], nil
	}
	return "", nil
}
