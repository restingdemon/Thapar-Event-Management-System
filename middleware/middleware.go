package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/restingdemon/thaparEvents/helpers"
)

var AuthenticationNotRequired map[string]bool = map[string]bool{
	"/create": true,
}

var RoleMethods = map[string][]string{
	"/users": {"superadmin"},
	// "/user-route":  "user",

}

// Authenticate is a middleware function that performs authentication
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		if AuthenticationNotRequired[requestedPath] {
			// If the requested path is in AuthenticationNotRequired, skip authentication
			next.ServeHTTP(w, r)
			return
		}

		clientToken := r.Header.Get("token")
		if clientToken == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("No auth header provided")))
			return
		}

		claims, msg := helpers.ValidateToken(clientToken)
		if msg != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token: " + msg))
			return
		}

		// Extract user roles from claims
		userType := claims.User_type

		// Check if any of the user's roles are authorized to access the requested route
		requiredRoles, ok := RoleMethods[requestedPath]
		if !ok {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Roles not defined for route: %s", requestedPath)))
			return
		}

		authorized := false
		for _, requiredRole := range requiredRoles {
			if strings.Contains(userType, requiredRole) {
				authorized = true
				break
			}
		}

		if !authorized {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("Access forbidden for route: %s", requestedPath)))
			return
		}

		ctx := context.WithValue(r.Context(), "email", claims.Email)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// contains checks if a string slice contains a specific value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
