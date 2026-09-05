package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxId int

const (
	requestID          = "X-Request-ID"
	requestCtxID ctxId = iota
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestID)
		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Add(requestID, id)
		ctx := context.WithValue(r.Context(), requestCtxID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RetrieveCTXIDFromContext(ctx context.Context) string {
	requestID := ctx.Value(requestCtxID).(string)
	return requestID
}
