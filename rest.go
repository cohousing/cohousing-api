package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"reflect"
	"strconv"
)

func ConfigureBasicTenantEndpoint(router *gin.RouterGroup, path string, domain interface{}) {
	tenantEndpoint := router.Group(path, MustBeTenant())

	repository := CreateRepository(reflect.TypeOf(domain))

	tenantEndpoint.GET("/", getResourceList(repository))
	tenantEndpoint.GET("/:id", getResourceById(repository))
	tenantEndpoint.POST("/", createNewResource(repository))
	tenantEndpoint.PUT("/:id", updateResource(repository))
	tenantEndpoint.DELETE("/:id", deleteResource(repository))
}

func getResourceList(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		list := repository.GetList(GetTenantFromContext(c))
		c.JSON(http.StatusOK, list)
	}
}

func getResourceById(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			if item, err := repository.GetById(GetTenantFromContext(c), id); err == nil {
				c.JSON(http.StatusOK, item)
			} else if err == gorm.ErrRecordNotFound {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "id is not an unsigned integer",
			})
		}
	}
}

func createNewResource(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func updateResource(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func deleteResource(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
