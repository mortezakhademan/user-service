package componentsJwt

import (
	"errors"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	_defaultRefreshTokenExpireTime = time.Hour * 24 * 365
	_defaultAccessTokenExpireTime  = time.Hour * 24
)

type Option func(*Jwt)

type rawObjectID []interface{}

type Jwt struct {
	issuer                     string
	refreshTokenExpireDuration time.Duration
	accessTokenExpireDuration  time.Duration
}

type JwtData struct {
	jwt.StandardClaims
	SessionID      *bson.ObjectID         `bson:"sess_id" json:"sess_id"`
	UserID         bson.ObjectID          `bson:"u_id" json:"u_id"`
	OrganizationID *bson.ObjectID         `json:"o_id,omitempty"`
	Permissions    map[int32][]int32      `bson:"perms" json:"perms,omitempty"`
	IsPrivateToken bool                   `bson:"is_p" json:"is_p,omitempty"`
	ValidatedToken bool                   `bson:"v_t" json:"v_t,omitempty"`
	ExtraData      map[string]interface{} `bson:"e_data" json:"e_data,omitempty"`
	IsRefreshToken bool                   `bson:"is_r" json:"is_r,omitempty"`
}

// NewJwt create new object from Jwt
func NewJwt(options ...Option) *Jwt {
	jwt := &Jwt{}
	jwt.refreshTokenExpireDuration = _defaultRefreshTokenExpireTime
	jwt.accessTokenExpireDuration = _defaultAccessTokenExpireTime
	for _, o := range options {
		o(jwt)
	}
	return jwt
}

// WithIssuer set issuer for jwt claim
func WithIssuer(issuer string) Option {
	return func(j *Jwt) {
		j.issuer = issuer
	}
}

// WithRefreshTokenExpireDuration set refresh token expired time
func WithRefreshTokenExpireDuration(duration time.Duration) Option {
	return func(j *Jwt) {
		j.refreshTokenExpireDuration = duration
	}
}

// WithAccessTokenExpireDuration set access token expired time
func WithAccessTokenExpireDuration(duration time.Duration) Option {
	return func(j *Jwt) {
		j.accessTokenExpireDuration = duration
	}
}

// CreateRefreshToken create new refresh token
func (j *Jwt) CreateRefreshToken(jwtData *JwtData, secretKey string) (string, error) {
	jwtData.ExpiresAt = time.Now().Add(j.refreshTokenExpireDuration).Unix()
	jwtData.IsRefreshToken = true
	return j.createToken(jwtData, secretKey)
}

// CreateAccessToken create new access token
func (j *Jwt) CreateAccessToken(jwtData *JwtData, secretKey string) (string, error) {
	jwtData.ExpiresAt = time.Now().Add(j.accessTokenExpireDuration).Unix()
	return j.createToken(jwtData, secretKey)
}

// ParseJwtToken parse jwt token payload
func ParseJwtToken(jwtToken string, secretKey string) (*JwtData, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, &jwtErr) {
			switch jwtErr.Errors {
			case jwt.ValidationErrorMalformed:
			case jwt.ValidationErrorUnverifiable:
			case jwt.ValidationErrorSignatureInvalid:
			case jwt.ValidationErrorAudience:
			case jwt.ValidationErrorExpired:
				return nil, componentsError.NewError(componentsError.ERROR_HTTP_TOKEN_EXPIRED)
			case jwt.ValidationErrorIssuedAt:
			case jwt.ValidationErrorIssuer:
			case jwt.ValidationErrorNotValidYet:
			case jwt.ValidationErrorId:
			case jwt.ValidationErrorClaimsInvalid:
			}
		}
		return nil, componentsError.NewError(componentsError.ERROR_HTTP_TOKEN_INVALID_SEGMENTS)
	}
	if !token.Valid {
		return nil, componentsError.NewError(componentsError.ERROR_HTTP_TOKEN_NOT_VALID)
	}
	data, err := parseJwtTokenClaims(claims)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Issuer get jwt issuer
func (j *Jwt) Issuer() string {
	return j.issuer
}

// IsTokenExpired check token is expired
func (j *JwtData) IsTokenExpired() bool {
	return j.ExpiresAt < time.Now().Unix()
}

// GetUserID get user id
func (j *JwtData) GetUserID() bson.ObjectID {
	return j.UserID
}

// GetSessionID get user session id
func (j *JwtData) GetSessionID() *bson.ObjectID {
	return j.SessionID
}

// GetServicePermissions get specific permissions code of an service
func (j *JwtData) GetServicePermissions(serviceCode int32) []int32 {
	return j.Permissions[serviceCode]
}

func (j *Jwt) createToken(jwtData *JwtData, secretKey string) (string, error) {
	if len(j.issuer) != 0 {
		jwtData.Issuer = j.issuer
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	strToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return strToken, nil

}

func parseJwtTokenClaims(claims jwt.MapClaims) (*JwtData, error) {
	data := &JwtData{}
	_, ok := claims["u_id"].(string)
	if !ok {
		claims["u_id"] = rawObjectIDToObjectID(claims["u_id"].([]interface{}))
	}
	_, ok = claims["sess_id"].(string)
	if !ok {
		claims["sess_id"] = rawObjectIDToObjectID(claims["sess_id"].([]interface{}))
	}
	b, err := bson.Marshal(claims)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(b, data); err != nil {
		return nil, err
	}
	return data, nil
}

func ParseJwtUnverified(token string) (*JwtData, error) {
	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(token, claims)
	if err != nil {
		return nil, err
	}
	return parseJwtTokenClaims(claims)
}

// ParseJwtToken parse jwt token payload
func IsValidJwtToken(jwtToken string, secretKey string) (bool, error) {
	token, err := new(jwt.Parser).Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, &jwtErr) {
			switch jwtErr.Errors {
			case jwt.ValidationErrorMalformed:
			case jwt.ValidationErrorUnverifiable:
			case jwt.ValidationErrorSignatureInvalid:
			case jwt.ValidationErrorAudience:
			case jwt.ValidationErrorExpired:
				return false, componentsError.NewError(componentsError.ERROR_HTTP_TOKEN_EXPIRED)
			case jwt.ValidationErrorIssuedAt:
			case jwt.ValidationErrorIssuer:
			case jwt.ValidationErrorNotValidYet:
			case jwt.ValidationErrorId:
			case jwt.ValidationErrorClaimsInvalid:
			}
		}
		return false, componentsError.NewError(componentsError.ERROR_HTTP_TOKEN_INVALID_SEGMENTS)
	}
	if !token.Valid {
		return false, componentsError.NewError(componentsError.ERROR_HTTP_TOKEN_NOT_VALID)
	}
	return true, nil
}

func rawObjectIDToObjectID(data rawObjectID) bson.ObjectID {
	userID := [12]uint8{}
	for k, datum := range data {
		if k == 12 {
			break
		}
		userID[k] = uint8(datum.(float64))
	}
	return userID
}
