package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/service/user/permission"
)

type wrappedResponseWriter struct {
	Rw     http.ResponseWriter
	Ignore bool
}

func (wrw *wrappedResponseWriter) WriteHeader(status int) {
	if status == 404 {
		wrw.Ignore = true
	} else {
		wrw.Rw.WriteHeader(status)
	}
}

func (wrw *wrappedResponseWriter) Write(p []byte) (int, error) {
	if wrw.Ignore {
		return 0, nil
	}
	return wrw.Rw.Write(p)
}

func (wrw *wrappedResponseWriter) Header() http.Header {
	return wrw.Rw.Header()
}

type wrappedHandler struct {
	h http.Handler
}

func (wh *wrappedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrw := wrappedResponseWriter{w, false}
	wh.h.ServeHTTP(&wrw, r)
	if wrw.Ignore == true {
		wrw.Rw.Header().Del("Content-Encoding")
		wrw.Rw.Header().Del("Vary")
		wrw.Rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		NotFoundHandler(wrw.Rw, r)
	}
}

func wrapHandler(handler http.Handler) http.Handler {
	return &wrappedHandler{handler}
}

// Make sure the user is a moderator, otherwise return forbidden
// TODO Clean this
func wrapModHandler(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := getUser(r)
		if userPermission.HasAdmin(currentUser) {
			handler(w, r)
		} else {
			http.Error(w, "admins only", http.StatusForbidden)
		}
	}
}
