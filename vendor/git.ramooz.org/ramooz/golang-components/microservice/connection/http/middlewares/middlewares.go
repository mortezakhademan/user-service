package middlewares

import (
	"context"
	componentsError "git.ramooz.org/ramooz/golang-components/error-handler"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

func ErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	s, ok := status.FromError(err)
	if !ok {
		runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
	}
	errCode := componentsError.GetHttpHeaderCode(int(s.Code()))
	if errCode != int(s.Code()) {
		customError := &runtime.HTTPStatusError{
			HTTPStatus: errCode,
			Err:        s.Err(),
		}
		runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, customError)
	} else {
		//if len(strconv.Itoa(errCode)) > 2 {
		//	w.WriteHeader(errCode)
		//}
		runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, s.Err())
	}
}

// AllowCORS add cors to http handler
func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization", "key"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}
