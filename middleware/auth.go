package middlewares

import (
	"context"
	"net/http"

	"github.com/hendraprawira/Tripatra-procurement/service"
)

type authString string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		bearer := "Bearer "
		auth = auth[len(bearer):]
		validate, err := service.JwtValidate(context.Background(), auth)
		if err != nil || !validate.Valid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// customClaim, _ := validate.Claims.(service.JwtCustomClaim)
		ctx := context.WithValue(r.Context(), authString("auth"), validate.Claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CtxValue(ctx context.Context) *service.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
	return raw
}
