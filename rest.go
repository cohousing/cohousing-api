package main

import (
	"fmt"
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
			if object, err := repository.GetById(GetTenantFromContext(c), id); err == nil {
				c.JSON(http.StatusOK, object)
			} else if err == gorm.ErrRecordNotFound {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
		} else {
			abortOnIdParsingError(c, id)
		}
	}
}

func createNewResource(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := reflect.New(repository.DomainType).Interface()
		if c.BindJSON(&object) == nil {

			createdObject, err := repository.Create(GetTenantFromContext(c), object)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			} else {
				c.JSON(http.StatusCreated, createdObject)
			}
		}
	}
}

func updateResource(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := reflect.New(repository.DomainType).Interface()
		if c.BindJSON(&object) == nil {
			if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
				if objectId := GetFieldByName(object, "ID").Uint(); objectId == id {
					updatedObject, err := repository.Update(GetTenantFromContext(c), object)

					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
					} else {
						c.JSON(http.StatusOK, updatedObject)
					}
				} else {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf("Id on path is different from id in object: %v != %v", id, objectId),
					})
				}
			} else {
				abortOnIdParsingError(c, id)
			}
		}
	}
}

func deleteResource(repository *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			if err = repository.Delete(GetTenantFromContext(c), id); err == nil {
				c.Status(http.StatusNoContent)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
		} else {
			abortOnIdParsingError(c, id)
		}
	}
}

func abortOnIdParsingError(c *gin.Context, id uint64) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": fmt.Sprintf("Id is not an unsigned integer: %v", id),
	})
}
