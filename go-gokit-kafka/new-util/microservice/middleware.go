package microservice

import (
	"context"
	"strconv"

	"git.bluebird.id/bluebird/util/uuid"

	jwt "github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//ErrUnauthorized is error for unauthorized access
	ErrUnauthorized = status.Error(codes.PermissionDenied, "Unauthorized access")
)

//AuthenticateMiddleware adds function to validate token
func AuthenticateMiddleware(signKey []byte, signMethod string) endpoint.Middleware {
	signing := jwt.GetSigningMethod(signMethod)
	keyFunc := func(*jwt.Token) (interface{}, error) {
		var sign interface{}
		var err error
		switch signing {
		case jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
			sign, err = jwt.ParseECPublicKeyFromPEM(signKey)
		case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512:
			sign, err = jwt.ParseRSAPublicKeyFromPEM(signKey)
		default:
			sign = signKey
		}
		if err != nil {
			return nil, err
		}
		return sign, nil
	}

	return endpoint.Chain(kitjwt.NewParser(keyFunc, signing, kitjwt.MapClaimsFactory), userContextMiddleware())
}

//AuthorizeMiddleware adds function to validate authorization
func AuthorizeMiddleware(serviceID, operationID int64) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			claims := ctx.Value(kitjwt.JWTClaimsContextKey)
			mapClaims := claims.(jwt.MapClaims)
			mapACL := mapClaims[CtxACL.String()]
			if mapACL == nil {
				return nil, ErrUnauthorized
			}
			acl := mapACL.(map[string]interface{})
			sid := strconv.Itoa(int(serviceID))
			access := acl[sid]
			if access == nil {
				return nil, ErrUnauthorized
			}

			oaccess := uint64(access.(float64))
			var one uint64 = 1
			oid := one << uint64(operationID-1)
			if oid&oaccess == 0 {
				return nil, ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

//AuthorizeExistMiddleware adds function to validate authorization at least one
func AuthorizeExistMiddleware(serviceID int64, operationID []int64) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			claims := ctx.Value(kitjwt.JWTClaimsContextKey)
			mapClaims := claims.(jwt.MapClaims)
			mapACL := mapClaims[CtxACL.String()]
			if mapACL == nil {
				return nil, ErrUnauthorized
			}
			acl := mapACL.(map[string]interface{})
			sid := strconv.Itoa(int(serviceID))
			access, ok := acl[sid]
			if !ok {
				return nil, ErrUnauthorized
			}

			oaccess := uint64(access.(float64))
			exist := false
			for _, opID := range operationID {
				var one uint64 = 1
				oid := one << uint64(opID-1)
				if oid&oaccess != 0 {
					exist = true
					break
				}
			}
			if !exist {
				return nil, ErrUnauthorized
			}

			return next(ctx, request)
		}
	}
}

//userContextMiddleware extracted user info from jwt claim
func userContextMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			claims := ctx.Value(kitjwt.JWTClaimsContextKey)
			mapClaims := claims.(jwt.MapClaims)
			return next(getUserContext(ctx, mapClaims), request)
		}
	}
}

func getUserContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	userContext := ctx
	uid, ok := claims[CtxUserID.String()]
	if ok {
		userID := int64(uid.(float64))
		userContext = context.WithValue(userContext, CtxUserID, userID)
	}

	domain, ok := claims[CtxDomain.String()]
	if ok {
		userContext = context.WithValue(userContext, CtxDomain, domain)
	}

	email, ok := claims[CtxEmail.String()]
	if ok {
		userContext = context.WithValue(userContext, CtxEmail, email)
	}

	phone, ok := claims[CtxPhone.String()]
	if ok {
		userContext = context.WithValue(userContext, CtxPhone, phone)
	}

	uidMsb, ok1 := claims[CtxUserIDMsb.String()]
	uidLsb, ok2 := claims[CtxUserIDLsb.String()]
	if ok1 && ok2 {
		userID := uuid.FromInt(uint64(uidMsb.(float64)), uint64(uidLsb.(float64)))
		userContext = context.WithValue(userContext, CtxUserUUID, userID)
	}

	return userContext
}
