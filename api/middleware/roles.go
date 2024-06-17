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
func isAutherizedToAccessResource(r *http.Request) (string, bool) {
	
	jwt := r.Header.Get("authorization")
	if jwt == "" {
		return "no jwt present", false
	}

	//get the token from the header
	token := strings.Split(jwt, "Bearer ")
	if len(token) != 2 {
		return "invalid token", false
	}

	return "unimplemented", false
}