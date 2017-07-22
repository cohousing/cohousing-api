package utils

import (
	"fmt"
	"github.com/cohousing/cohousing-api/config"
	"github.com/cohousing/cohousing-api/domain/admin"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
)

const (
	GIN_TENANT    = "gin_tenant"
	GIN_IS_TENANT = "gin_is_tenant"
)

// Resolves the context based on URL
func ContextResolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host

		host, err := trimHost(host)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		tenant := config.GetTenantByHost(host)
		if tenant != nil {
			c.Set(GIN_TENANT, tenant)
			c.Set(GIN_IS_TENANT, true)
		} else {
			c.Set(GIN_IS_TENANT, false)
		}

		c.Next()
	}
}

func trimHost(host string) (string, error) {
	if strings.Index(host, ":") > -1 {
		var err error
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return "", err
		} else {
			return host, nil
		}
	} else {
		return host, nil
	}
}

func GetTenantFromContext(c *gin.Context) *admin.Tenant {
	return c.MustGet(GIN_TENANT).(*admin.Tenant)
}

func IsTenantRequest(c *gin.Context) bool {
	return c.GetBool(GIN_IS_TENANT)
}

func IsNotTenantRequest(c *gin.Context) bool {
	return !IsTenantRequest(c)
}

func MustBeTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsNotTenantRequest(c) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("No tenant found on URL: %s", c.Request.Host),
			})
		}
	}
}

func MustBeAdminDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := location.Get(c)
		host, err := trimHost(url.Host)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			if host != config.GetConfig().AdminDomain {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": fmt.Sprintf("Endpoint requested (%s) not found on domain: %s", c.Request.RequestURI, url.Host),
				})
			}
		}

	}
}
