package middlewareAuth

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwtAuth "github.com/ahsansandiah/dpo-test/packages/auth/jwt"
	"github.com/ahsansandiah/dpo-test/packages/config"
	jsonResponse "github.com/ahsansandiah/dpo-test/packages/json"
	logger "github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/google/uuid"
)

type Middleware interface {
	InitLog(next http.Handler) http.Handler
	CheckToken(next http.Handler) http.Handler
	GetTokenInHeader(r *http.Request) (string, error)
}

type Options struct {
	jwt       jwtAuth.Jwt
	secretKey string
	log       logger.Log
	json      jsonResponse.Json
}

func NewMiddleware(cfg *config.Config, lg logger.Log, jsonRes jsonResponse.Json) Middleware {
	opt := new(Options)
	opt.jwt = jwtAuth.NewJwt(cfg)
	opt.secretKey = cfg.JwtSecretKey
	opt.log = lg
	opt.json = jsonRes

	return opt
}

func (o *Options) InitLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set body
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// set constants
		tnow := time.Now()
		id := uuid.New().String()

		ctx := context.WithValue(r.Context(), config.ContextKey("body"), bodyBytes)
		ctx = context.WithValue(ctx, config.ContextKey("startTime"), tnow)
		ctx = context.WithValue(ctx, config.ContextKey("id"), id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (o *Options) GetTokenInHeader(r *http.Request) (string, error) {
	authzHeader := r.Header.Get("Authorization")
	if authzHeader == "" {
		return "", ErrorAuthHeaderEmpty
	}

	accessToken := strings.Split(authzHeader, " ")
	if accessToken[0] != "Bearer" {
		err := ErrorAuthNotHaveBearer
		return "", err
	}

	if len(accessToken) == 1 {
		err := ErrorAuthNotHaveToken
		return "", err
	}

	return accessToken[1], nil
}

func (o *Options) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accessToken, err := o.GetTokenInHeader(r)
		if err != nil {
			o.json.ErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		_, err = o.jwt.ExtractJwtToken(accessToken)
		if err != nil {
			o.json.ErrorResponse(w, r, http.StatusUnauthorized, ErrorInvalidTokenOrExpired)
			return
		}

		jwtData, err := o.jwt.VerifyAccessToken(accessToken, o.secretKey)
		if err != nil {
			o.json.ErrorResponse(w, r, http.StatusUnauthorized, ErrorInvalidTokenOrExpired)
			return
		}

		if jwtData == nil {
			o.json.ErrorResponse(w, r, http.StatusUnauthorized, ErrorAccessTokenEmpty)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
