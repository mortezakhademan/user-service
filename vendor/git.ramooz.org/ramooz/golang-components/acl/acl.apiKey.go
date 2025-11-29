package acl

import (
	"context"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
)

// NewAclFromApiKey create acl object for control api key
func NewAclFromApiKey(ctx context.Context, serviceCode int32, apiKey string, options ...Option) (*Acl, error) {
	return newAclFromApiKey(ctx, serviceCode, apiKey, applyOption(options...))
}

func newAclFromApiKey(ctx context.Context, serviceCode int32, apiKey string, config *Config) (*Acl, error) {
	if config.getApiKeyInfoFunc == nil {
		return nil, componentsError.NewError(ERROR_API_KEY_FUNC_INVALID)
	}
	userId, userPermissions, err := config.getApiKeyInfoFunc(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	config.validateAcl = false
	acl, err := newAcl(serviceCode, *userId, nil, userPermissions, nil, config)
	if err != nil {
		return nil, err
	}
	acl.validatedPermissions = true
	return acl, nil
}
