package scsheader

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
)

type HeaderSession struct {
	*scs.SessionManager
}

func (s *HeaderSession) LoadAndSaveHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := "Authorization"

		token := r.Header.Get(headerKey)
		if token != "" {
			token = strings.TrimPrefix(token, "Token ")
		}
		ctx, err := s.Load(r.Context(), token)
		if err != nil {
			_ = log.Output(2, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		bw := &bufferedResponseWriter{ResponseWriter: w}
		sr := r.WithContext(ctx)
		next.ServeHTTP(bw, sr)

		if s.Status(ctx) == scs.Modified {
			token, _, err := s.Commit(ctx)
			if err != nil {
				_ = log.Output(2, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.Header().Set(headerKey, "Token "+token)
		}

		if bw.code != 0 {
			w.WriteHeader(bw.code)
		}
		_, _ = w.Write(bw.buf.Bytes())
	})
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf  bytes.Buffer
	code int
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	bw.code = code
}
