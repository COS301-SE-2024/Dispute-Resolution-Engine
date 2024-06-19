package middleware

import (
	"api/models"
	"api/utilities"
	"net/http"
	// "github.com/gorilla/mux"
)

type Role struct {
	roles map[int][]string
}

// type RestrictedRoleEndpointFunction func(router *mux.Router, h Handler)

// type RestrictedRoleEndpoint struct {
// 	accessLevels map[int][]string
// 	accessLevel  int
// 	inner        Handler
// }

// NewRole creates a new Role struct
func NewRole() *Role {
	accessLevels := map[int][]string{
		0: {"guest"},
		1: {"user"},
		2: {"mediator"},
		3: {"adjudicator"},
		4: {"arbitrator"},
		5: {"systemAnalyst"},
		6: {"systemAdmin"},
	}
	return &Role{roles: accessLevels}
}

// func NewRestrictedRoleEndpoint(inner Handler, accessLevel int) *RestrictedRoleEndpoint {
// 	return &RestrictedRoleEndpoint{
// 		accessLevels: NewRole().roles,
// 		accessLevel:  accessLevel,
// 		inner:        inner,
// 	}
// }

func (r *Role) matchKeyToValue(value string) (int, bool) {
	for key, values := range r.roles {
		for _, v := range values {
			if v == value {
				return key, true
			}
		}
	}
	return -1, false
}

func RoleMiddleware(next http.Handler) http.Handler {
	roles := NewRole()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil {
			utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		//string match the role to the map
		userRole := claims.User.Role

		//check if the role is in the map
		_, ok := roles.matchKeyToValue(userRole)

		if !ok {
			utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}

		//check if the role is allowed to access the resource

		next.ServeHTTP(w, r)
	})
}

// Role struct to store role data
// func isAuthorizedToAccessResource(r *http.Request, authLevel int) (string, bool) {

// 	claims := GetClaims(r)
// 	if claims == nil {
// 		return "jwt error", false
// 	}

// 	role := claims.User.Role
// 	accessLevels := NewRole().roles

// }