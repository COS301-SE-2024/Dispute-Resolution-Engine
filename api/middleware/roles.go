package middleware

import (
	"net/http"
	"strings"
)

type Role struct {
	roles map[int][]string
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



// Role struct to store role data
// func isAuthorizedToAccessResource(r *http.Request, authLevel int) (string, bool) {
	
// 	claims := GetClaims(r)
// 	if claims == nil {
// 		return "jwt error", false
// 	}

// 	role := claims.User.Role
// 	accessLevels := NewRole().roles


// }