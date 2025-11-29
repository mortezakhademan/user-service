package acl

import (
	"context"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	componentsJwt "git.ramooz.org/ramooz/golang-components/jwt"
)

type AclJwt struct {
	*Acl
	JwtData *componentsJwt.JwtData
	token   string
}

// NewWithJwt create acl object for control jwt token
func NewWithJwt(serviceCode int32, jwtToken string, options ...Option) (*AclJwt, error) {
	return newWithJwt(serviceCode, jwtToken, applyOption(options...))

}

func newWithJwt(serviceCode int32, jwtToken string, config *Config) (*AclJwt, error) {
	aclJwt := &AclJwt{
		token: jwtToken,
	}
	jwtData, err := componentsJwt.ParseJwtUnverified(jwtToken)
	if err != nil {
		return nil, componentsError.New(ERROR_JWT_TOKEN_IS_INVALID, []string{"error in pars jwt"})
	}

	aclJwt.JwtData = jwtData
	secretKey := config.publicSecretKey
	if jwtData.IsPrivateToken {
		secretKey = config.privateSecretKey
	}

	if len(secretKey) == 0 {
		return nil, componentsError.NewError(ERROR_SECRET_KEY_IS_EMPTY)
	}

	if valid, err := componentsJwt.IsValidJwtToken(jwtToken, secretKey); !valid {
		return nil, err
	}
	if jwtData.ValidatedToken {
		config.validateAcl = false
	}
	if acl, err := newAcl(serviceCode, jwtData.UserID, jwtData.SessionID, jwtData.Permissions, nil, config); err != nil {
		return nil, err
	} else {
		aclJwt.Acl = acl
	}
	if jwtData.ValidatedToken {
		aclJwt.validatedPermissions = true
	}
	return aclJwt, nil
}

// GetAclJwtFromContext get acl jwt from context
func GetAclJwtFromContext(ctx context.Context) (*AclJwt, error) {
	if jwt, ok := ctx.Value(_defaultAclContextKey).(*AclJwt); ok {
		return jwt, nil
	}
	return nil, componentsError.NewError(ERR_EXTRACT_ACL_FROM_CONTEXT)
}

// SetAclToContext set acl jwt into context
func (a *AclJwt) SetAclToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, _defaultAclContextKey, a)
}

// GetExtraData get extra data from jwtData
func (a *AclJwt) GetExtraData() map[string]interface{} {
	return a.JwtData.ExtraData
}
