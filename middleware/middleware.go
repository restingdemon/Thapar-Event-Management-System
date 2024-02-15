package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/utils"
)

var AuthenticationNotRequired map[string]bool = map[string]bool{
	"/create": true,

	"/events/get": true, 

	"/soc/get":true,

}

var RoleMethods = map[string][]string{
	"/users/get":       {utils.AdminRole, utils.SuperAdminRole},
	"/users/update/":   {utils.AdminRole, utils.UserRole, utils.SuperAdminRole},
	"/soc/register":    {utils.SuperAdminRole},
	"/soc/update/":     {utils.AdminRole, utils.SuperAdminRole},
	"/event/create":    {utils.AdminRole, utils.SuperAdminRole},
	"/event/update/":   {utils.AdminRole, utils.SuperAdminRole},
	"/event/register/": {utils.UserRole},
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

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("No auth header provided")))
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token type"))
			return
		}

		clientToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, msg := helpers.ValidateToken(clientToken)
		if msg != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token: " + msg))
			return
		}

		// Extract user roles from claims
		userType := claims.User_type
		userEmail := claims.Email
		//userId := claims.Id
		// Check if any of the user's roles are authorized to access the requested route
		authorized := false
		for path, requiredRoles := range RoleMethods {
			if strings.HasPrefix(requestedPath, path) {
				for _, requiredRole := range requiredRoles {
					if strings.Contains(userType, requiredRole) {
						authorized = true
						break
					}
				}
				break
			}
		}

		if !authorized {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("Access forbidden for route: %s", requestedPath)))
			return
		}
		ctx, err := CheckHTTPAuthorization(r, r.Context(), userType, userEmail)
		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("Permission denied: %s", err)))
			return
		}
		sjhda := ctx.Value("email")
		fmt.Println(sjhda)

		// Call the next handler in the chain with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
