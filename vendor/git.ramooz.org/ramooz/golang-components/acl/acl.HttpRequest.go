package acl

import (
	"fmt"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	"net/http"
	"strconv"
)

// NewAclFromHttpRequest create acl from http request base on api key or jwt token in context value
func NewAclFromHttpRequest(request *http.Request, serviceCode int32, options ...Option) (AclController, error) {
	config := applyOption(options...)
	token, err := GetBearerTokenFromHttpRequest(request)
	if err == nil {
		jwtAcl, err := newWithJwt(serviceCode, token, config)
		if err != nil {
			return nil, err
		}
		return jwtAcl, nil
	}

	if errDetail, ok := err.(*componentsError.Error); ok {
		_, httpCode, customCode := errDetail.GetSplitCode()
		errCode, _ := strconv.Atoi(fmt.Sprintf("%d%d", httpCode, customCode))
		if errCode == ERROR_NO_HEADER_IN_REQUEST {
			apiKey, err := GetApiKeyFromHttpRequest(request, config.apiKeyHttpHeaderName)
			if err == nil {
				return newAclFromApiKey(request.Context(), serviceCode, apiKey, config)
			}
		}
	}
	return nil, err
}

func GetAclContextFromHttp(req *http.Request, serviceCode int32, validateWithUserService bool, jwtPublicSecret string, jwtPrivateSecret string) (*AclContext, error) {
	aclOpts := aclOptions(validateWithUserService, jwtPublicSecret, jwtPrivateSecret)

	aclController, err := NewAclFromHttpRequest(req, serviceCode, aclOpts...)
	if err != nil {
		return nil, err
	}
	aclCtx := &AclContext{
		Context:       req.Context(),
		AclController: aclController,
	}
	return aclCtx, nil
}

func aclOptions(validate bool, publicSecret, privateSecret string) []Option {
	opts := []Option{}

	if validate {
		opts = append(opts, WithValidateACL(true))
	}

	if len(publicSecret) != 0 {
		opts = append(opts, WithPublicSecretKey(publicSecret))
	}

	if len(privateSecret) != 0 {
		opts = append(opts, WithPrivateSecretKey(privateSecret))
	}

	return opts
}
