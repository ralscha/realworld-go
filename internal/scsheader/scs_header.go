package scsheader

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
)

const headerKey string = "Authorization"

type HeaderSession struct {
	*scs.SessionManager
}

func (s *HeaderSession) LoadAndSaveHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var sr *http.Request
		var ctx context.Context

		token := r.Header.Get(headerKey)
		if token != "" {
			token = strings.TrimPrefix(token, "Token ")
		}

		var err error
		ctx, err = s.Load(r.Context(), token)
		if err != nil {
			_ = log.Output(2, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sr = r.WithContext(ctx)

		next.ServeHTTP(w, sr)
	})
}
