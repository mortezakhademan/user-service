package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"git.ramooz.org/ramooz/golang-components/acl"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	pbCache "git.ramooz.org/ramooz/pb/apis-gen/imports/cache/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jhump/protoreflect/desc"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
	"time"

	errors "git.ramooz.org/ramooz/golang-components/microservice/error"
	"git.ramooz.org/ramooz/golang-components/microservice/helper"
	"git.ramooz.org/ramooz/golang-components/microservice/serviceInfo"
)

const (
	REQUIRED = iota
	OPTIONAL
	WITHOUT_CREDENTIAL
)

// PermissionDescriptor get permissions and credentialType from descriptor, default validate is true in implemention
//
//	credential type:
//	  * REQUIRED = 0;
//	  * OPTIONAL = 1;
//	  * WITHOUT_CREDENTIAL = 2;
//
// returns serviceCode, method permissions,credentialType, validate
type PermissionDescriptor func(ctx context.Context) (int32, []int32, int32, bool)

type validator interface {
	ValidateAll() error
}

type validatorLegacy interface {
	Validate() error
}

func Middlewares(middlewares ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(
		middlewares...,
	)
}

// MethodDescriptors save methods descriptors into context for any methods
func MethodDescriptors(descriptors map[string]*desc.MethodDescriptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md := descriptors[info.FullMethod]
		ctx = metadata.AppendToOutgoingContext(context.WithValue(ctx, "desc", md))
		return handler(ctx, req)
	}
}

type ExtractCacheDescriptor func(context.Context, string) ([]pbCache.CacheKey, int32, []int32, []int32)

func GrpcGatewayCache(cacheDescriptor ExtractCacheDescriptor, publicSecretKey string, privateSecretKey string, serviceCode int32) grpc.UnaryClientInterceptor {
	c := cache.New(5*time.Minute, 10*time.Minute)
	aclOpts := aclOptions(nil, nil, false, "", publicSecretKey, privateSecretKey)
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var httpMethod string
		if md, ok := metadata.FromOutgoingContext(ctx); !ok {
			return componentsError.New(errors.ErrorHttpMethodNotFound, nil)
		} else {
			httpMethod = md.Get("http-method")[0]
			if httpMethod != "GET" {
				return invoker(ctx, method, req, reply, cc, opts...)
			}
		}
		cacheKeys, ttl, acceptPermissions, excludePermissions := cacheDescriptor(ctx, method)
		if ttl == -1 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		aclController, _ := acl.NewAclFromOutgoingContext(ctx, serviceCode, aclOpts...)
		if aclController == nil && (len(acceptPermissions) > 0 || len(excludePermissions) > 0) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		if len(acceptPermissions) > 0 {
			if !aclController.HasAnyPermissionsAccess(acceptPermissions...) {
				return invoker(ctx, method, req, reply, cc, opts...)
			}
		}
		if len(excludePermissions) > 0 {
			if aclController.HasAnyPermissionsAccess(excludePermissions...) {
				return invoker(ctx, method, req, reply, cc, opts...)
			}
		}
		key := method + "***" + fmt.Sprint(req)
		for i := 0; i < len(cacheKeys); i++ {
			switch cacheKeys[i] {
			case pbCache.CacheKey_CK_HeaderAuthorize:
				if aclController == nil {
					return invoker(ctx, method, req, reply, cc, opts...)
				}
				key += aclController.GetUserID().Hex()
			}
		}
		if data, ok := c.Get(key); ok {
			if jdata, err := json.Marshal(data); err != nil {
				return err
			} else {
				if err := json.Unmarshal(jdata, reply); err != nil {
					return err
				}
			}
			return nil
		}
		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
			return err
		}
		c.Set(key, reply, time.Second*time.Duration(ttl))
		return nil
	}
}

// GrpcAuth middleware for check jwt and api key with permissions
func GrpcAuth(permissionDescriptor PermissionDescriptor, GetUserPermissionsFunc acl.GetUserPermissions, GetApiKeyInfoFunc acl.GetApiKeyInfo, apiKeyCustomContextKeyName, jwtPublicSecret, jwtPrivateSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		serviceCode, permissions, credentialType, validate := permissionDescriptor(ctx)

		if credentialType == WITHOUT_CREDENTIAL {
			if token, _ := acl.GetBearerTokenFromGrpcIncomingContext(ctx); token == "" {
				return handler(ctx, req)
			}
			return nil, componentsError.New(componentsError.ERROR_HTTP_FORBIDDEN, []string{"users with auth token can't access to this api"})
		}

		aclOpts := aclOptions(GetUserPermissionsFunc, GetApiKeyInfoFunc, validate, apiKeyCustomContextKeyName, jwtPublicSecret, jwtPrivateSecret)

		aclController, err := acl.NewAclFromIncomingContext(ctx, serviceCode, aclOpts...)
		if err != nil {
			if credentialType == OPTIONAL && isErrorOptional(err) {
				return handler(ctx, req)
			}
			return nil, err
		}

		if len(permissions) != 0 {
			if !aclController.HasAnyPermissionsAccessInOtherService(serviceCode, permissions...) {
				return nil, componentsError.New(componentsError.ERROR_HTTP_FORBIDDEN, []string{fmt.Sprintf("don't have access to service code:%d don't have any of this permissions: %+v", serviceCode, permissions)})
			}
		}

		ctx = aclController.SetAclToContext(ctx)
		return handler(ctx, req)
	}
}

// GrpcJwt middleware for check jwt
func GrpcJwt(permissionDescriptor PermissionDescriptor, serviceInfo *serviceInfo.ServiceInfo, GetUserPermissionsFunc acl.GetUserPermissions, jwtPublicSecret, jwtPrivateSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		serviceCode, permissions, credentialType, validate := permissionDescriptor(ctx)

		if credentialType == WITHOUT_CREDENTIAL {
			if jwtToken, _ := acl.GetBearerTokenFromGrpcIncomingContext(ctx); jwtToken == "" {
				return handler(ctx, req)
			}
			return nil, componentsError.New(componentsError.ERROR_HTTP_FORBIDDEN, []string{"users with auth token can't access to this api"})
		}

		jwtToken, err := acl.GetBearerTokenFromGrpcIncomingContext(ctx)
		if err != nil {
			if credentialType == OPTIONAL && isErrorOptional(err) {
				return handler(ctx, req)
			}
			return nil, err
		}

		aclOpts := aclOptions(GetUserPermissionsFunc, nil, validate, "", jwtPublicSecret, jwtPrivateSecret)

		aclController, err := acl.NewWithJwt(serviceInfo.Code, jwtToken, aclOpts...)
		if err != nil {
			return nil, err
		}

		if len(permissions) != 0 {
			if !aclController.HasAnyPermissionsAccessInOtherService(serviceCode, permissions...) {
				return nil, componentsError.New(componentsError.ERROR_HTTP_FORBIDDEN, []string{fmt.Sprintf("don't have access to service code:%d don't have any of this permissions: %+v", serviceCode, permissions)})
			}
		}

		ctx = aclController.SetAclToContext(ctx)
		return handler(ctx, req)

	}
}

// GrpcAPIKey middleware for check api key with permissions
func GrpcAPIKey(permissionDescriptor PermissionDescriptor, serviceInfo *serviceInfo.ServiceInfo, GetApiKeyInfoFunc acl.GetApiKeyInfo, apiKeyCustomContextKeyName, jwtPublicSecret, jwtPrivateSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		serviceCode, permissions, credentialType, validate := permissionDescriptor(ctx)

		if credentialType == WITHOUT_CREDENTIAL {
			if apiKey, _ := acl.GetApiKeyFromContext(ctx, apiKeyCustomContextKeyName); apiKey == "" {
				return handler(ctx, req)
			}

			return nil, componentsError.New(componentsError.ERROR_HTTP_FORBIDDEN, []string{"users with auth token can't access to this api"})
		}

		apiKey, err := acl.GetApiKeyFromContext(ctx, apiKeyCustomContextKeyName)
		if err != nil {
			if credentialType == OPTIONAL && isErrorOptional(err) {
				return handler(ctx, req)
			}
			return nil, err
		}

		aclOpts := aclOptions(nil, GetApiKeyInfoFunc, validate, apiKeyCustomContextKeyName, jwtPublicSecret, jwtPrivateSecret)

		aclController, err := acl.NewAclFromApiKey(ctx, serviceInfo.Code, apiKey, aclOpts...)
		if err != nil {
			return nil, err
		}

		if len(permissions) != 0 {
			if !aclController.HasAnyPermissionsAccessInOtherService(serviceCode, permissions...) {
				return nil, componentsError.New(componentsError.ERROR_HTTP_FORBIDDEN, []string{fmt.Sprintf("don't have access to service code:%d don't have any of this permissions: %+v", serviceCode, permissions)})
			}
		}

		ctx = aclController.SetAclToContext(ctx)
		return handler(ctx, req)
	}
}

// GrpcRecovery recovery panics
func GrpcRecovery() grpc.UnaryServerInterceptor {
	rec := func(p interface{}) (err error) {
		err = status.Errorf(codes.Unknown, "%v", p)
		helper.Log.Errorf("panic triggered: Error %v led to gRPC server recovery", err)
		return
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(rec),
	}
	return grpc_recovery.UnaryServerInterceptor(opts...)
}

// GrpcValidator validate your message fields, for user validator please check https://github.com/envoyproxy/protoc-gen-validate
func GrpcValidator() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		switch in := req.(type) {
		case validator:
			if err := in.ValidateAll(); err != nil {
				return nil, componentsError.New(errors.ERR_VALIDATE_FIELDS, []string{err.Error()})
			}
		case validatorLegacy:
			if err := in.Validate(); err != nil {
				return nil, componentsError.New(errors.ERR_VALIDATE_FIELDS, []string{err.Error()})
			}
		}
		return handler(ctx, req)
	}
}

// ClientInterceptor show client requested
func ClientInterceptor() grpc.DialOption {
	client := func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		start := time.Now()
		if err := invoker(ctx, method, req, reply, cc, opts...); err == nil {
			helper.Log.Infof("Invoked RPC method=%s; Duration=%s", method, time.Since(start))
		}
		return nil
	}
	return grpc.WithUnaryInterceptor(client)
}

func aclOptions(GetUserPermissionsFunc acl.GetUserPermissions, GetApiKeyInfoFunc acl.GetApiKeyInfo, validate bool, apiKeyCustomContextKeyName, publicSecret, privateSecret string) []acl.Option {
	opts := []acl.Option{}

	if validate {
		opts = append(opts, acl.WithValidateACL(true))
	}

	if len(publicSecret) != 0 {
		opts = append(opts, acl.WithPublicSecretKey(publicSecret))
	}

	if len(privateSecret) != 0 {
		opts = append(opts, acl.WithPrivateSecretKey(privateSecret))
	}

	if GetUserPermissionsFunc != nil {
		opts = append(opts, acl.WithGetUserPermissionsFunction(GetUserPermissionsFunc))
	}

	if GetApiKeyInfoFunc != nil {
		opts = append(opts, acl.WithGetAPIKeyInfoFunction(GetApiKeyInfoFunc))
	}

	if len(apiKeyCustomContextKeyName) != 0 {
		opts = append(opts, acl.WithAPIKeyCustomContextName(apiKeyCustomContextKeyName))
	}

	return opts
}

func isErrorOptional(err error) bool {
	if cErr, ok := err.(*componentsError.Error); ok {
		_, httpCode, customCode := cErr.GetSplitCode()
		errCode, _ := strconv.Atoi(fmt.Sprintf("%d%d", httpCode, customCode))
		if errCode == acl.ERROR_NO_HEADER_IN_REQUEST || errCode == acl.ERROR_AUTHORIZATION_NOT_FOUND {
			return true
		}
	}
	return false
}
