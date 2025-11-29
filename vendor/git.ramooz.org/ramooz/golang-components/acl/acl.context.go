package acl

import (
	"context"
	"errors"
	"fmt"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	"strconv"
)

const (
	_defaultAccessTokenContextKey = "authorization"
	_defaultApiKeyContextKey      = "api_key"
	_defaultAclContextKey         = "acl"
	_defaultServiceContextKey     = "s_code"
)

type AclContext struct {
	context.Context
	AclController
}

func NewAclFromOutgoingContext(ctx context.Context, serviceCode int32, options ...Option) (AclController, error) {
	config := applyOption(options...)
	token, err := GetBearerTokenFromOutgoingContext(ctx)
	if err != nil {
		return handleGetBearerTokenFromContextError(ctx, config, serviceCode, err)
	}
	return NewAclByToken(token, config, serviceCode)
}

// NewAclFromIncomingContext create acl from context base on api key or jwt token in context value
func NewAclFromIncomingContext(ctx context.Context, serviceCode int32, options ...Option) (AclController, error) {
	config := applyOption(options...)
	token, err := GetBearerTokenFromGrpcIncomingContext(ctx)
	if err != nil {
		return handleGetBearerTokenFromContextError(ctx, config, serviceCode, err)
	}
	return NewAclByToken(token, config, serviceCode)
}

func handleGetBearerTokenFromContextError(ctx context.Context, config *Config, serviceCode int32, err error) (AclController, error) {
	var errDetail *componentsError.Error
	if errors.As(err, &errDetail) {
		_, httpCode, customCode := errDetail.GetSplitCode()
		errCode, _ := strconv.Atoi(fmt.Sprintf("%d%d", httpCode, customCode))
		if errCode == ERROR_NO_HEADER_IN_REQUEST {
			apiKey, err := GetApiKeyFromContext(ctx, config.apiKeyContextName)
			if err == nil {
				return newAclFromApiKey(ctx, serviceCode, apiKey, config)
			}
		}
	}
	return nil, err
}

func NewAclByToken(token string, config *Config, serviceCode int32) (AclController, error) {
	jwtAcl, err := newWithJwt(serviceCode, token, config)
	if err != nil {
		return nil, err
	}
	return jwtAcl, nil
}

func GetAclContext(ctx context.Context) (AclContext, error) {
	acl, err := GetAclFromContext(ctx)
	if err != nil {
		aclJwt, err := GetAclJwtFromContext(ctx)
		if err != nil {
			return AclContext{nil, acl}, err
		}
		return AclContext{ctx, aclJwt}, nil
	}
	return AclContext{ctx, acl}, nil
}
