package middleware

import (
	"api/models"
	"api/utilities"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
	roles map[int][]string
}

// type RestrictedRoleEndpointFunction func(router *mux.Router, h Handler)

// type RestrictedRoleEndpoint struct {
// 	accessLevels map[int][]string
// 	accessLevel  int
// 	inner        Handler
// }

var roleMiddleware Role
var onceRole sync.Once

func NewRole() Role {
	onceRole.Do(func() {
		roleMiddleware = createRole()
	})
	return roleMiddleware
}

// NewRole creates a new Role struct
func createRole() Role {
	accessLevels := map[int][]string{
		0: {"guest"},
		1: {"user"},
		2: {"mediator"},
		3: {"adjudicator"},
		4: {"arbitrator"},
		5: {"systemAnalyst"},
		6: {"systemAdmin"},
	}
	return &RoleMiddleware{roles: accessLevels}
}

// func NewRestrictedRoleEndpoint(inner Handler, accessLevel int) *RestrictedRoleEndpoint {
// 	return &RestrictedRoleEndpoint{
// 		accessLevels: NewRole().roles,
// 		accessLevel:  accessLevel,
// 		inner:        inner,
// 	}
// }

func (r *RoleMiddleware) matchKeyToValue(value string) (int, bool) {
	for key, values := range r.roles {
		for _, v := range values {
			if v == value {
				return key, true
			}
		}
	}
	return -1, false
}

func (r *RoleMiddleware) RoleMiddleware(reqAuthlevel int) gin.HandlerFunc {
	logger := utilities.NewLogger().LogWithCaller()
	jwt := NewJwtMiddleware()

	return func(c *gin.Context) {
		claims := jwt.GetClaims(c)
		if claims == nil {
			logger.Error("No claims")
			c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		//string match the role to the map
		userRole := claims.User.Role

		//check if the role is in the map
		authLevel, ok := r.matchKeyToValue(userRole)

		if !ok {
			logger.Error("Role not found")
			c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}

		//check if the role is allowed to access the resource
		if authLevel < reqAuthlevel {
			logger.Error("Unauthorized")
			c.JSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
        c.Next()
	}
}
