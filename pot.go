package pot

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/getfrontierhq/buf-plugins/binder"
	potErrors "github.com/getfrontierhq/buf-plugins/errors"
	"github.com/go-chi/chi/v5"
)

type (
	HandlerFunc       func(ctx context.Context, req interface{}) (interface{}, error)
	MiddlewareFunc    func(HandlerFunc) HandlerFunc
	DecoderFunc       func(req interface{}) error
	MethodHandlerFunc func(ctx context.Context, srv interface{}, dec DecoderFunc, middleware MiddlewareFunc) (interface{}, error)

	MethodDescriptor struct {
		MethodName string

		HttpMethod string
		HttpPath   string
		Handler    MethodHandlerFunc
	}

	ServiceDescriptor struct {
		ServiceName string
		HandlerType interface{}
		Methods     []MethodDescriptor
	}
)

func NewDecoderFunc(r *http.Request) DecoderFunc {
	dec := binder.RequestDecoder{Request: r}

	return func(req interface{}) error {
		return dec.Bind(req)
	}
}

func httpHandlerWrapper(impl interface{}, handler MethodHandlerFunc) http.HandlerFunc {
	type ErrResp struct {
		Message string `json:"message"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		decoder := NewDecoderFunc(r)
		out, err := handler(r.Context(), impl, decoder, nil)
		if err == nil {
			encoder := &binder.ResponseEncoder{ResponseWriter: rw}
			err = encoder.BindBody(out)
			if err == nil {
				return
			}
		}

		statusCode, err := potErrors.ParseErr(err)

		potErr := &potErrors.Error{}
		if errors.As(err, &potErr) {
			rw.WriteHeader(statusCode)
			json.NewEncoder(rw).Encode(err)
			return
		}

		rw.WriteHeader(statusCode)
		json.NewEncoder(rw).Encode(ErrResp{Message: err.Error()})
	}
}

func RegisterServiceWithChi(desc *ServiceDescriptor, impl interface{}, router chi.Router) http.Handler {
	if impl != nil {
		ht := reflect.TypeOf(desc.HandlerType).Elem()
		st := reflect.TypeOf(impl)
		if !st.Implements(ht) {
			log.Fatalf("pot: RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}

	for _, method := range desc.Methods {
		switch method.HttpMethod {
		case http.MethodGet:
			router.Get(method.HttpPath, httpHandlerWrapper(impl, method.Handler))
		case http.MethodPost:
			router.Post(method.HttpPath, httpHandlerWrapper(impl, method.Handler))
		case http.MethodPut:
			router.Put(method.HttpPath, httpHandlerWrapper(impl, method.Handler))
		case http.MethodPatch:
			router.Patch(method.HttpPath, httpHandlerWrapper(impl, method.Handler))
		case http.MethodDelete:
			router.Delete(method.HttpPath, httpHandlerWrapper(impl, method.Handler))
		default:
			panic("pot: RegisterService found unsupported HTTP method: " + method.HttpMethod)
		}
	}

	return router
}

func RegisterService(desc *ServiceDescriptor, impl interface{}) http.Handler {
	router := chi.NewRouter()
	return RegisterServiceWithChi(desc, impl, router)
}
