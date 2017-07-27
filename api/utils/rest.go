package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	RecordsPerPage int = 50
)

type LinkFactory func(halResource domain.HalResource, basePath string, detailed bool)

type BasicEndpointConfig struct {
	Path            string
	Domain          interface{}
	domainType      reflect.Type
	LinkFactory     LinkFactory
	DBFactory       db.DBFactory
	RouterHandlers  []gin.HandlerFunc
	GetListHandlers []gin.HandlerFunc
	GetHandlers     []gin.HandlerFunc
	CreateHandlers  []gin.HandlerFunc
	UpdateHandlers  []gin.HandlerFunc
	DeleteHandlers  []gin.HandlerFunc
}

func ConfigureBasicEndpoint(router *gin.RouterGroup, config BasicEndpointConfig) *gin.RouterGroup {
	config.domainType = reflect.TypeOf(config.Domain)

	endpoint := router.Group(config.Path, config.RouterHandlers...)

	repository := db.CreateRepository(config.domainType, config.DBFactory)

	endpoint.GET("/", append(config.GetListHandlers, getResourceList(config, endpoint.BasePath(), repository))...)
	endpoint.GET("/:id", append(config.GetHandlers, getResourceById(config, endpoint.BasePath(), repository))...)
	endpoint.POST("/", append(config.CreateHandlers, createNewResource(config, endpoint.BasePath(), repository))...)
	endpoint.PUT("/:id", append(config.UpdateHandlers, updateResource(config, endpoint.BasePath(), repository))...)
	endpoint.DELETE("/:id", append(config.DeleteHandlers, deleteResource(repository))...)

	return endpoint
}

func getResourceList(config BasicEndpointConfig, basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		lookupObject, page := parseQuery(c, config.domainType)
		list, count := repository.GetList(c, lookupObject, GetStartRecord(page, RecordsPerPage), RecordsPerPage)
		valueList := reflect.ValueOf(list).Elem()
		listLength := valueList.Len()
		domainList := CreateObjectList(basePath, make([]interface{}, listLength), page, count, RecordsPerPage)
		for i := 0; i < listLength; i++ {
			object := valueList.Index(i).Addr().Interface()
			domainList.Objects[i] = object
			addLinks(object, config.LinkFactory, basePath, false)
		}

		c.JSON(http.StatusOK, domainList)
	}
}

func getResourceById(config BasicEndpointConfig, basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			if object, err := repository.GetById(c, id); err == nil {
				addLinks(object, config.LinkFactory, basePath, true)
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

func createNewResource(config BasicEndpointConfig, basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := reflect.New(config.domainType).Interface()
		if err := c.ShouldBindWith(&object, binding.JSON); err == nil {

			createdObject, err := repository.Create(c, object)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			} else {
				addLinks(createdObject, config.LinkFactory, basePath, true)
				c.JSON(http.StatusCreated, createdObject)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
	}
}

func updateResource(config BasicEndpointConfig, basePath string, repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := reflect.New(config.domainType).Interface()
		if c.BindJSON(&object) == nil {
			if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
				if objectId := GetFieldByName(object, "ID").Uint(); objectId == id {
					updatedObject, err := repository.Update(c, object)

					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
					} else {
						addLinks(updatedObject, config.LinkFactory, basePath, true)
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

func deleteResource(repository *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			if err = repository.Delete(c, id); err == nil {
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

func addLinks(object interface{}, linkFactory LinkFactory, basePath string, detailed bool) {
	halResource, ok := object.(domain.HalResource)
	if ok {
		linkFactory(halResource, basePath, detailed)
	}
}

func abortOnIdParsingError(c *gin.Context, id uint64) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": fmt.Sprintf("Id is not an unsigned integer: %v", id),
	})
}

func parseQuery(c *gin.Context, domainType reflect.Type) (lookupObject interface{}, pageNumber int) {
	page := GetCurrentPage(c)

	var buffer bytes.Buffer
	buffer.WriteString("{")
	queryParams := c.Request.URL.Query()
	for queryParam := range queryParams {
		if queryParam != "page" {
			if buffer.Len() > 1 {
				buffer.WriteString(",")
			}
			buffer.WriteString("\"")
			buffer.WriteString(strings.Replace(queryParam, "\"", "\\\"", -1))
			buffer.WriteString("\":")

			queryValue := queryParams.Get(queryParam)
			_, err := strconv.ParseInt(queryValue, 10, 64)
			if err == nil {
				buffer.WriteString(queryValue)
			} else {
				buffer.WriteString("\"")
				buffer.WriteString(strings.Replace(queryValue, "\"", "\\\"", -1))
				buffer.WriteString("\"")
			}
		}
	}
	buffer.WriteString("}")

	lookupObject = reflect.New(domainType).Interface()
	err := json.Unmarshal(buffer.Bytes(), lookupObject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
	}

	return lookupObject, page
}
