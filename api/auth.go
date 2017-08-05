package api

import (
	"fmt"
	"github.com/cohousing/cohousing-tenant-api/config"
	"github.com/cohousing/cohousing-tenant-api/db"
	"github.com/cohousing/cohousing-tenant-api/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type AuthOperation string

const (
	GIN_TOKEN                       = "gin_token"
	GIN_CLAIMS                      = "gin_token"
	GIN_AUTHENTICATED               = "gin_authenticated"
	PERM_CREATE       AuthOperation = "Create"
	PERM_READ         AuthOperation = "Read"
	PERM_UPDATE       AuthOperation = "Update"
	PERM_DELETE       AuthOperation = "Delete"
)

type CohousingClaims struct {
	jwt.StandardClaims
	Permissions domain.Permission `json:"permissions,omitempty"`
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authorizationHeader, "Bearer ") {
			authorizationValue := authorizationHeader[7:]

			var claims CohousingClaims
			token, err := jwt.ParseWithClaims(authorizationValue, &claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(config.GetConfig().TokenSecret), nil
			})

			if token.Valid {
				c.Set(GIN_AUTHENTICATED, true)
				c.Set(GIN_TOKEN, token)
				c.Set(GIN_CLAIMS, claims)
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "That's not even a token",
					})
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "Token expired",
					})
				} else {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf("Couldn't handle this token: %v", err),
					})
				}
			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Couldn't handle this token: %v", err),
				})
			}
		}
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "You must provide a username",
		})
		return
	}

	password := c.PostForm("password")
	if password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "You must provide a password",
		})
		return
	}

	redirect := c.PostForm("redirect_url")
	if redirect == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "You must provide a redirect_url",
		})
		return
	}
	redirectUrl, err := url.Parse(redirect)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Redirect URL couldn't be parse: %v", err),
		})
		return
	}

	var user domain.User
	db.GetTenantDB(GetTenantFromContext(c)).Where("`username` = ?", username).First(&user)
	db.GetTenantDB(GetTenantFromContext(c)).Model(&user).Related(&user.Groups, "Groups")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	permissions := user.ResolvePermissions()
	claims := CohousingClaims{}
	claims.Subject = fmt.Sprintf("%d", user.ID)
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
	claims.Permissions = permissions
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConfig().TokenSecret))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Couldn't sign token: %v", err))
		return
	}

	query := redirectUrl.Query()
	query.Add("id_token", tokenString)
	redirectUrl.RawQuery = query.Encode()
	c.Redirect(http.StatusFound, redirectUrl.String())
}

func MustAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsAuthenticated(c) == false {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Next()
		}
	}
}

func AuthorizeDomainObject(domain interface{}, operation AuthOperation) gin.HandlerFunc {
	domainType := reflect.TypeOf(domain).Name()
	permission := fmt.Sprintf("%s%ss", operation, domainType)

	return func(c *gin.Context) {
		resolvedPermissions := GetPermissionsFromContext(c)
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
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func MustBeGlobalAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		resolvedPermissions := GetPermissionsFromContext(c)

		if resolvedPermissions.GlobalAdmin {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func IsAuthenticated(c *gin.Context) bool {
	return c.GetBool(GIN_AUTHENTICATED)
}

func GetPermissionsFromContext(c *gin.Context) domain.Permission {
	claims := c.MustGet(GIN_CLAIMS).(CohousingClaims)
	return claims.Permissions
}
