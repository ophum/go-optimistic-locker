package optimisticlocker

import (
	"log"
	"net/http"

	versionmanager "github.com/ophum/go-optimistic-locker/version_manager"
)

type VersionKeyParser func(r *http.Request) (string, error)

type Middleware = func(next http.Handler) http.Handler

type Locker interface {
	PreconditionCheck(keyMaker VersionKeyParser) Middleware

	VersionManager() versionmanager.VersionManager
}

type locker struct {
	versionManager versionmanager.VersionManager
}

func NewLocker(vm versionmanager.VersionManager) Locker {
	return &locker{vm}
}

func (l *locker) PreconditionCheck(keyParser VersionKeyParser) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ifMatch := r.Header.Get("If-Match")
			if ifMatch == "" {
				w.WriteHeader(http.StatusPreconditionRequired)
				return
			}
			key, err := keyParser(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				return
			}
			version, err := l.versionManager.Get(r.Context(), key)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			if ifMatch != version {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (l *locker) VersionManager() versionmanager.VersionManager {
	return l.versionManager
}
