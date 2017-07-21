package api

import (
	"fmt"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"net/http"
	"reflect"
	"strconv"
)

func ConfigureBasicTenantEndpoint(router *gin.RouterGroup, path string, domain interface{}) *gin.RouterGroup {
	tenantEndpoint := router.Group(path, MustBeTenant())

	repository := db.CreateRepository(reflect.TypeOf(domain))

	tenantEndpoint.GET("/", getResourceList(tenantEndpoint.BasePath(), repository))
	tenantEndpoint.GET("/:id", getResourceById(tenantEndpoint.BasePath(), repository))
	tenantEndpoint.POST("/", createNewResource(tenantEndpoint.BasePath(), repository))
	tenantEndpoint.PUT("/:id", updateResource(tenantEndpoint.BasePath(), repository))
	tenantEndpoint.DELETE("/:id", deleteResource(tenantEndpoint.BasePath(), repository))

	return tenantEndpoint
}

func getResourceList(basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		list := repository.GetList(GetTenantFromContext(c))
		valueList := reflect.ValueOf(list).Elem()
		listLength := valueList.Len()
		domainList := &ObjectList{
			Objects: make([]interface{}, listLength),
		}
		for i := 0; i < listLength; i++ {
			object := valueList.Index(i).Addr().Interface()
			domainList.Objects[i] = object
			halResource, ok := object.(domain.HalResource)
			if ok {
				halResource.Populate(false)
			}
		}
		addPaginationLinks(domainList, basePath, 1, 1, 2)

		c.JSON(http.StatusOK, domainList)
	}
}

func getResourceById(basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			if object, err := repository.GetById(GetTenantFromContext(c), id); err == nil {
				halResource, ok := object.(domain.HalResource)
				if ok {
					halResource.Populate(true)
				}
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

func createNewResource(basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := reflect.New(repository.DomainType).Interface()
		if err := c.ShouldBindWith(&object, binding.JSON); err == nil {

			createdObject, err := repository.Create(GetTenantFromContext(c), object)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			} else {
				c.JSON(http.StatusCreated, createdObject)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
	}
}

func updateResource(basePath string, repository *db.Repository) gin.HandlerFunc {
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

func deleteResource(basePath string, repository *db.Repository) gin.HandlerFunc {
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
