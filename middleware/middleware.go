package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/utils"
)

var AuthenticationNotRequired map[string]bool = map[string]bool{
	"/create":         true,
	"/event/get":      true,
	"/soc/get":        true,
	"/soc/get/events": true,
	"/feedback":       true,
}

var RoleMethods = map[string][]string{
	"/users/get":               {utils.UserRole, utils.AdminRole, utils.SuperAdminRole},
	"/users/update/":           {utils.AdminRole, utils.UserRole, utils.SuperAdminRole},
	"/users/get/registrations": {utils.UserRole},
	"/soc/register":            {utils.SuperAdminRole},
	"/soc/update/":             {utils.AdminRole, utils.SuperAdminRole},
	"/soc/get/notvisible":      {utils.SuperAdminRole},
	"/soc/get/allevents":       {utils.AdminRole, utils.SuperAdminRole},
	//"/event/create":               {utils.AdminRole, utils.SuperAdminRole},
	"/event/update/":              {utils.AdminRole, utils.SuperAdminRole},
	"/event/register/":            {utils.UserRole},
	"/event/get/registrations/":   {utils.AdminRole, utils.SuperAdminRole},
	"/event/check/registrations/": {utils.UserRole},
	"/event/visibility/":          {utils.SuperAdminRole},
	"/event/delete/":              {utils.AdminRole, utils.SuperAdminRole},
	"/event/upload/":              {utils.AdminRole},
	"/event/photo/delete":         {utils.AdminRole},
	"/event/poster/upload":        {utils.AdminRole},
	"/soc/dashboard":              {utils.AdminRole, utils.SuperAdminRole},
	"/event/dashboard":            {utils.AdminRole, utils.SuperAdminRole},
	"/event/get/notvisible":       {utils.SuperAdminRole},
}

// Authenticate is a middleware function that performs authentication
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		if AuthenticationNotRequired[requestedPath] || strings.HasPrefix(requestedPath, "/event/create") {

			if strings.HasPrefix(requestedPath, "/event/create") {
				apiKey := r.Header.Get("X-API-Key")
				if apiKey == "" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("No API key provided"))
					return
				} else if apiKey == "integration" {
					next.ServeHTTP(w, r)
					return
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Not Authorized"))
					return
				}
			}
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
