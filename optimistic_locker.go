package optimisticlocker

import (
	"net/http"
)

type ResourceGetter func(r *http.Request) (any, error)

type Middleware = func(next http.Handler) http.Handler

func PreconditionCheck(resourceGetter ResourceGetter, opts ...PreconditionCheckOption) Middleware {
	option := PreconditionCheckOptions{
		EtagGenerator: DefaultEtagGenerator,
	}
	for _, opt := range opts {
		opt.apply(&option)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ifMatch := r.Header.Get("If-Match")
			if ifMatch == "" {
				w.WriteHeader(http.StatusPreconditionRequired)
				return
			}
			targetResource, err := resourceGetter(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			etag, err := option.EtagGenerator(targetResource)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if ifMatch != etag {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
