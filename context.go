package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

const (
	GIN_TENANT    = "gin_tenant"
	GIN_IS_TENANT = "gin_is_tenant"
)

// Resolves the context based on URL
func ContextResolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		host, _, _ := net.SplitHostPort(c.Request.Host)

		tenant := GetTenantByHost(host)
		if tenant != nil {
			c.Set(GIN_TENANT, tenant)
			c.Set(GIN_IS_TENANT, true)
		} else {
			c.Set(GIN_IS_TENANT, false)
		}

		c.Next()
	}
}

func GetTenantFromContext(c *gin.Context) *Tenant {
	return c.MustGet(GIN_TENANT).(*Tenant)
}

func MustBeTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool(GIN_IS_TENANT) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("No tenant found on URL: %s", c.Request.Host),
			})
		}
	}
}
