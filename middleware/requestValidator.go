package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/utils"
)

func CheckHTTPAuthorization(r *http.Request, ctx context.Context, userType string, userEmail string) (context.Context, error) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/users/update"):
		// extracting email from path parameters
		vars := mux.Vars(r)
		email, ok := vars["email"]
		if !ok {
			return ctx, fmt.Errorf("no email provided")
		}
		if userType == "superadmin" {
			ctx = context.WithValue(ctx, "email", email)
			return ctx, nil
		}
		if email != userEmail {
			return ctx, fmt.Errorf("you can only update your own details")
		}

		ctx = context.WithValue(ctx, "email", email)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/soc/register"):
		if userType != utils.SuperAdminRole {
			return ctx, fmt.Errorf("Invalid Role")
		}
	case strings.HasPrefix(r.URL.Path, "/soc/update"):
		vars := mux.Vars(r)
		email, ok := vars["email"]
		if !ok {
			return ctx, fmt.Errorf("no email provided")
		}
		if userType == utils.SuperAdminRole {
			ctx = context.WithValue(ctx, "email", email)
			ctx = context.WithValue(ctx, "role", utils.SuperAdminRole)
			return ctx, nil
		} else if userType == utils.AdminRole {
			if email != userEmail {
				return ctx, fmt.Errorf("you can only update your own details")
			}
			ctx = context.WithValue(ctx, "email", userEmail)
			ctx = context.WithValue(ctx, "role", userType)
			return ctx, nil
		} else {
			return ctx, fmt.Errorf("Invalid Role")
		}
	case strings.HasPrefix(r.URL.Path, "/event/create"):
		vars := mux.Vars(r)
		email, ok := vars["email"]
		if !ok {
			return ctx, fmt.Errorf("no email provided")
		}
		if userType == "superadmin" {
			ctx = context.WithValue(ctx, "email", email)
			return ctx, nil
		}
		if email != userEmail {
			return ctx, fmt.Errorf("you can only create event in your own society")
		}

		ctx = context.WithValue(ctx, "email", email)
		return ctx, nil
	case strings.HasPrefix(r.URL.Path, "/event/update"):
		vars := mux.Vars(r)
		eventId, ok := vars["eventId"]
		if !ok {
			return ctx, fmt.Errorf("no Event Id provided")
		}
		if userType == utils.SuperAdminRole {
			ctx = context.WithValue(ctx, "eventId", eventId)
			ctx = context.WithValue(ctx, "role", utils.SuperAdminRole)
			return ctx, nil
		} else if userType == utils.AdminRole {
			// if email != userEmail {
			// 	return ctx, fmt.Errorf("you can only update your own details")
			// }
			ctx = context.WithValue(ctx, "eventId", eventId)
			ctx = context.WithValue(ctx, "role", userType)
			ctx = context.WithValue(ctx, "email", userEmail)
			return ctx, nil
		} else {
			return ctx, fmt.Errorf("Invalid Role")
		}
	}

	// Default to allowing access if the route is not explicitly handled
	return ctx, nil
}
