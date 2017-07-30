package api

import (
	"encoding/base64"
	"fmt"
	"github.com/cohousing/cohousing-tenant-api/db"
	"github.com/cohousing/cohousing-tenant-api/domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type AuthOperation string

const (
	GIN_USER                        = "gin_user"
	GIN_AUTHENTICATED               = "gin_authenticated"
	PERM_CREATE       AuthOperation = "Create"
	PERM_READ         AuthOperation = "Read"
	PERM_UPDATE       AuthOperation = "Update"
	PERM_DELETE       AuthOperation = "Delete"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if len(authorizationHeader) == 0 {
			c.Set(GIN_AUTHENTICATED, false)
		} else if strings.HasPrefix(authorizationHeader, "Basic ") {
			authorizationValue := authorizationHeader[6:]
			authorization, err := base64.StdEncoding.DecodeString(authorizationValue)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't decode base64: %v\n", err)
				abortWithUnauthenticated(c)
				c.Set(GIN_AUTHENTICATED, false)
				return
			}

			userPassArray := strings.Split(string(authorization), ":")
			if len(userPassArray) != 2 {
				fmt.Fprintf(os.Stderr, "Basic auth not constructed correctly: %v\n", authorization)
				abortWithUnauthenticated(c)
				c.Set(GIN_AUTHENTICATED, false)
				return
			}

			username := userPassArray[0]
			password := userPassArray[1]

			var user domain.User
			db.GetTenantDB(GetTenantFromContext(c)).Where("`username` = ?", username).First(&user)
			db.GetTenantDB(GetTenantFromContext(c)).Model(&user).Related(&user.Groups, "Groups")

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Password not correct for request %s: %v\n", authorization, err)
				abortWithUnauthenticated(c)
				c.Set(GIN_AUTHENTICATED, false)
				return
			}

			c.Set(GIN_USER, &user)
			c.Set(GIN_AUTHENTICATED, true)
		}
	}
}

func MustAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsAuthenticated(c) == false {
			abortWithUnauthenticated(c)
		} else {
			c.Next()
		}
	}
}

func AuthorizeDomainObject(domain interface{}, operation AuthOperation) gin.HandlerFunc {
	domainType := reflect.TypeOf(domain).Name()
	permission := fmt.Sprintf("%s%ss", operation, domainType)

	return func(c *gin.Context) {
		u := GetUserFromContext(c)
		resolvedPermissions := u.ResolvePermissions()
		hasPermission, err := resolvedPermissions.HasPermission(permission)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if resolvedPermissions.GlobalAdmin || hasPermission {
			c.Next()
		} else {
			abortWithUnauthenticated(c)
		}
	}
}

func MustBeGlobalAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := GetUserFromContext(c)
		resolvedPermissions := u.ResolvePermissions()

		if resolvedPermissions.GlobalAdmin {
			c.Next()
		} else {
			abortWithUnauthenticated(c)
		}
	}
}

func IsAuthenticated(c *gin.Context) bool {
	return c.GetBool(GIN_AUTHENTICATED)
}

func GetUserFromContext(c *gin.Context) *domain.User {
	return c.MustGet(GIN_USER).(*domain.User)
}

func ResolvePermission(c *gin.Context) domain.Permission {
	return GetUserFromContext(c).ResolvePermissions()
}

func abortWithUnauthenticated(c *gin.Context) {
	// Credentials doesn't match, we return 401 and abort handlers chain.
	c.Header("WWW-Authenticate", "Basic realm=Cohousing")
	c.AbortWithStatus(401)
}
