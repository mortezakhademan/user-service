package acl

import (
	"context"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	componentsJwt "git.ramooz.org/ramooz/golang-components/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc/metadata"
	"strconv"
)

// AclController is acl methods for auth type
type AclController interface {
	HasAccess(permissionCode int32) bool
	HasAccessInOtherService(permissionCode int32, serviceCode int32) bool
	NotHasAccess(permissionCode int32) bool
	HasAnyPermissionsAccess(permissionCodes ...int32) bool
	HasAnyPermissionsAccessInOtherService(serviceCode int32, permissionCodes ...int32) bool
	HasAllPermissionsAccess(permissionCodes ...int32) bool
	GetPrivateToken(extraData map[string]any) (string, error)
	SetPrivateTokenToOutgoingContext(ctx context.Context, serviceCode int32, extraData map[string]any) (context.Context, error)
	SetAclToContext(ctx context.Context) context.Context
	GetUserID() bson.ObjectID
	GetSessionID() *bson.ObjectID
}

type Acl struct {
	userPermissions map[int32][]int32

	config *Config

	validatedPermissions bool
	privateJwtToken      string
	userId               bson.ObjectID
	sessionId            *bson.ObjectID
	extraData            map[string]interface{}
}

// New create new acl object
func New(serviceCode int32, userId bson.ObjectID, sessionId *bson.ObjectID, userPermissions map[int32][]int32, extraData map[string]interface{}, options ...Option) (*Acl, error) {
	return newAcl(serviceCode, userId, sessionId, userPermissions, extraData, applyOption(options...))
}

func newAcl(serviceCode int32, userId bson.ObjectID, sessionId *bson.ObjectID, userPermissions map[int32][]int32, extraData map[string]interface{}, config *Config) (*Acl, error) {
	if serviceCode == 0 {
		return nil, componentsError.NewError(ERROR_SERVICE_CODE_HAS_BEEN_NOT_SET)
	}
	config.currentServiceCode = serviceCode
	acl := &Acl{userPermissions: userPermissions, userId: userId, sessionId: sessionId, extraData: extraData, config: config}
	if config.validateAcl {
		if !acl.IsValidAcl() {
			return nil, componentsError.NewError(ERROR_ACL_DATA_NOT_VALID)
		}
	}
	return acl, nil
}

func GetAclFromContext(ctx context.Context) (*Acl, error) {
	acl, ok := ctx.Value(_defaultAclContextKey).(*Acl)
	if !ok {
		return nil, componentsError.NewError(ERR_EXTRACT_ACL_FROM_CONTEXT)
	}
	return acl, nil
}

// HasAccess check has access to service method
func (acl *Acl) HasAccess(permissionCode int32) bool {
	if acl == nil {
		return false
	}
	return acl.HasAccessInOtherService(acl.config.currentServiceCode, permissionCode)
}

// HasAnyPermissionsAccess check permissions has access to one permission (or)
func (acl *Acl) HasAnyPermissionsAccess(permissionCodes ...int32) bool {
	if acl == nil {
		return false
	}
	return acl.HasAnyPermissionsAccessInOtherService(acl.config.currentServiceCode, permissionCodes...)
}

func (acl *Acl) HasAnyPermissionsAccessInOtherService(serviceCode int32, permissionCodes ...int32) bool {
	if acl == nil {
		return false
	}
	for _, code := range permissionCodes {
		if acl.HasAccessInOtherService(serviceCode, code) {
			return true
		}
	}
	return false

}

// HasAllPermissionsAccess check permissions has access to all permissions (and)
func (acl *Acl) HasAllPermissionsAccess(permissionCodes ...int32) bool {
	if acl == nil {
		return false
	}
	for _, code := range permissionCodes {
		if !acl.HasAccess(code) {
			return false
		}
	}
	return true
}

// HasAccessInOtherService check has access to other service method
func (acl *Acl) HasAccessInOtherService(serviceCode int32, permissionCode int32) bool {
	if acl == nil {
		return false
	}

	for _, userPermissionCode := range acl.userPermissions[serviceCode] {
		if userPermissionCode == permissionCode {
			return true
		}
	}
	return false
}

// NotHasAccess check not has access to service method
func (acl *Acl) NotHasAccess(permissionCode int32) bool {
	if acl == nil {
		return true
	}

	return !acl.HasAccess(permissionCode)
}

// IsValidAcl validate permissions from user service
func (acl *Acl) IsValidAcl() bool {
	if acl == nil {
		return false
	}
	if acl.validatedPermissions {
		return true
	}
	if acl.config.getUserPermissionsFunc != nil {
		isValid := acl.isValidPermissions(acl.config.getUserPermissionsFunc(acl.userId))
		if isValid {
			acl.validatedPermissions = true
		}
		return isValid
	}
	return false
}

// GetPrivateToken get private jwt token from acl
// extraData append to old extraData
func (acl *Acl) GetPrivateToken(extraData map[string]any) (string, error) {
	if len(acl.config.privateSecretKey) == 0 {
		return "", componentsError.NewError(ERROR_PRIVATE_SECRET_KEY_IS_EMPTY)
	}
	if len(acl.privateJwtToken) != 0 {
		return acl.privateJwtToken, nil
	}
	for k, v := range acl.extraData {
		if _, ok := extraData[k]; !ok {
			extraData[k] = v
		}
	}
	privateJwtData := &componentsJwt.JwtData{
		IsPrivateToken: true,
		ValidatedToken: acl.validatedPermissions,
		Permissions:    acl.userPermissions,
		ExtraData:      extraData,
		SessionID:      acl.sessionId,
		UserID:         acl.userId,
	}
	jwtToken, err := componentsJwt.NewJwt().CreateAccessToken(privateJwtData, acl.config.privateSecretKey)
	if err != nil {
		return "", err
	}
	acl.privateJwtToken = jwtToken
	return jwtToken, nil
}

// SetPrivateTokenToOutgoingContext
// extraData append to acl extra data
func (acl *Acl) SetPrivateTokenToOutgoingContext(ctx context.Context, serviceCode int32, extraData map[string]any) (context.Context, error) {
	token, err := acl.GetPrivateToken(extraData)
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md.Set("Authorization", "Bearer "+token)
	md[_defaultServiceContextKey] = append(md[_defaultServiceContextKey], strconv.Itoa(int(serviceCode)))
	outCtx := metadata.NewOutgoingContext(ctx, md)
	return outCtx, nil
}

func (acl *Acl) SetAclToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, _defaultAclContextKey, acl)
}

// GetUserID return user id
func (acl *Acl) GetUserID() bson.ObjectID {
	if acl == nil {
		return bson.NilObjectID
	}
	return acl.userId
}

// GetSessionID return session id
func (acl *Acl) GetSessionID() *bson.ObjectID {
	if acl == nil {
		return nil
	}
	return acl.sessionId
}

func (acl *Acl) isValidPermissions(userPermissions map[int32][]int32) bool {
	if acl == nil {
		return false
	}
	for userServiceCode, userPerms := range userPermissions {
		for aclServiceCode, aclPerms := range acl.userPermissions {
			if userServiceCode == aclServiceCode {
				if !acl.isAclWithUserPermissionsMatch(aclPerms, userPerms) {
					return false
				}
			}
		}
	}
	return true
}

func (acl *Acl) isAclWithUserPermissionsMatch(aclPerms, userPerms []int32) bool {
	if len(aclPerms) != len(userPerms) {
		return false
	}
	for i := range aclPerms {
		if aclPerms[i] != userPerms[i] {
			return false
		}
	}
	return true
}

// GetExtraData get extra data from jwtData
func (acl *Acl) GetExtraData() map[string]interface{} {
	return acl.extraData
}
