package acl

import (
	"context"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

func GetBearerTokenFromHttpRequest(request *http.Request) (string, error) {
	authorization := request.Header.Get("Authorization")
	if authorization == "" {
		return "", componentsError.NewError(ERROR_AUTHORIZATION_NOT_FOUND)
	}
	token := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))
	return token, nil
}

func GetApiKeyFromHttpRequest(request *http.Request, apiKeyHeaderName string) (string, error) {
	if len(apiKeyHeaderName) == 0 {
		apiKeyHeaderName = _defaultApiKeyContextKey
	}
	return "", componentsError.New(ERR_METHOD_NOT_IMPLEMENTED, []string{"GetApiKeyFromHttpRequest not implemented"})
}

func GetBearerTokenFromOutgoingContext(ctx context.Context) (string, error) {
	foundedHeaders, err := extractHeaderFromOutgoingContext(ctx, _defaultAccessTokenContextKey)
	if err != nil {
		return "", err
	}
	if len(foundedHeaders) != 1 {
		return "", componentsError.NewError(ERROR_AUTHORIZATION_NOT_FOUND)
	}
	return getBearerTokenFromHeader(foundedHeaders[0])
}
func GetBearerTokenFromGrpcIncomingContext(ctx context.Context) (string, error) {
	foundedHeaders, err := extractHeaderFromIncomingContext(ctx, _defaultAccessTokenContextKey)
	if err != nil {
		return "", err
	}
	if len(foundedHeaders) != 1 {
		return "", componentsError.NewError(ERROR_AUTHORIZATION_NOT_FOUND)
	}
	return getBearerTokenFromHeader(foundedHeaders[0])
}
func getBearerTokenFromHeader(headerToken string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(headerToken, prefix) {
		return "", componentsError.NewError(ERROR_MISSING_BEARER_IN_HEADER)
	}
	token := strings.TrimPrefix(headerToken, prefix)
	return token, nil
}

func SetBearerTokenToOutgoingContext(ctx context.Context, bearerToken string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md.Set("Authorization", "Bearer "+bearerToken)
	outCtx := metadata.NewOutgoingContext(ctx, md)
	return outCtx, nil
}

func GetApiKeyFromContext(ctx context.Context, apiKeyHeaderName string) (string, error) {
	if len(apiKeyHeaderName) == 0 {
		apiKeyHeaderName = _defaultApiKeyContextKey
	}
	foundedHeaders, err := extractHeaderFromIncomingContext(ctx, apiKeyHeaderName)
	if err != nil {
		return "", err
	}
	if len(foundedHeaders) != 1 {
		return "", componentsError.NewError(ERROR_AUTHORIZATION_NOT_FOUND)
	}

	apiKey := foundedHeaders[0]
	return apiKey, nil
}

func extractHeaderFromMetaData(md metadata.MD, header string) ([]string, error) {
	foundedHeaders, ok := md[header]
	if !ok {
		return nil, componentsError.NewError(ERROR_NO_HEADER_IN_REQUEST)
	}

	return foundedHeaders, nil
}

func extractHeaderFromIncomingContext(ctx context.Context, header string) ([]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, componentsError.NewError(ERROR_NO_HEADER_IN_REQUEST)
	}
	return extractHeaderFromMetaData(md, header)
}

func extractHeaderFromOutgoingContext(ctx context.Context, header string) ([]string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, componentsError.NewError(ERROR_NO_HEADER_IN_REQUEST)
	}
	return extractHeaderFromMetaData(md, header)
}
