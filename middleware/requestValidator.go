package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func CheckHTTPAuthorization(r *http.Request, ctx context.Context, userType string, userEmail string) (context.Context, error) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/users/get"):
		if userType == "superadmin" {
			vars := mux.Vars(r)
			email, ok := vars["email"]
			if !ok {
				return ctx, fmt.Errorf("no email provided")
			}
			ctx = context.WithValue(ctx, "email", email)
			return ctx, nil
		}
		ctx = context.WithValue(ctx, "email", userEmail)
		sjhda := ctx.Value("email")
		fmt.Println(sjhda)
		return ctx, nil
	}

	// Default to allowing access if the route is not explicitly handled
	return ctx, nil
}
