package handlers

import (
	"expvar"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"net/http/pprof"
)

// SetRuntimeAsRootHandler set runtime mux as root handler http server mux
func SetRuntimeAsRootHandler(mux *http.ServeMux, rMux *runtime.ServeMux) *http.ServeMux {
	mux.Handle("/", rMux)
	return mux
}

// DebuggerHandler add pprof handlers handlers to http server mux
func DebuggerHandler(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())
	return mux
}

// SwaggerHandler add swagger file embedded to http handler path, swaggerFileName (swagger.json or swagger.yaml and etc)
func SwaggerHandler(mux *http.ServeMux, swaggerFileName string, swagger []byte) *http.ServeMux {
	mux.HandleFunc("/"+swaggerFileName, func(w http.ResponseWriter, _ *http.Request) {
		w.Write(swagger)
	})
	return mux
}
