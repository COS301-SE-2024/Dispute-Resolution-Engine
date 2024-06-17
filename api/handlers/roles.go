package handlers

import "github.com/gorilla/mux"

type Role struct {
	roles map[int][]string
}

type RestrictedRoleEndpointFunction func(router *mux.Router, h Handler)

type RestrictedRoleEndpoint struct {
	accessLevels map[int][]string
	accessLevel  int
	inner        Handler
}

// NewRole creates a new Role struct
func NewRole() *Role {
	accessLevels := map[int][]string{
		0: {"unverified"},
		1: {"user"},
		2: {"mediator"},
		3: {"adjudicator"},
		4: {"arbitrator"},
		5: {"overseer"},
		6: {"systemAnalyst"},
		7: {"systemAdmin"},
	}
	return &Role{roles: accessLevels}
}

func NewRestrictedRoleEndpoint(inner Handler, accessLevel int) *RestrictedRoleEndpoint {
	return &RestrictedRoleEndpoint{
		accessLevels: NewRole().roles,
		accessLevel:  accessLevel,
		inner:        inner,
	}
}

func NewRestrictedRoleEndpointFunction(inner Handler, accessLevel int) RestrictedRoleEndpointFunction {
	return func(router *mux.Router, h Handler) {
		
		
	}
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