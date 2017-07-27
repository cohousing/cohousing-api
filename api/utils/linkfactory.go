package utils

import (
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
	"reflect"
)

var (
	factory map[string]LinkFactory = make(map[string]LinkFactory)
)

type LinkFactory func(c *gin.Context, halResource domain.HalResource, basePath string, detailed bool)

func AddLinkFactory(resourceType interface{}, linkFactory LinkFactory) {
	factory[reflect.TypeOf(resourceType).Name()] = linkFactory
}

// Add links to resource based on the type of it
func AddLinks(c *gin.Context, resource interface{}, basePath string, detailed bool) {
	resourceType := reflect.TypeOf(resource).Elem()

	addLinks := func(halResource domain.HalResource) {
		linkFactory := factory[resourceType.Name()]
		linkFactory(c, halResource, basePath, detailed)
	}

	if resourceType.Kind() == reflect.Array {
		resourceList := resource.([]interface{})
		for i := 0; i < len(resourceList); i++ {
			addLinks(resourceList[i].(domain.HalResource))
		}
	} else {
		addLinks(resource.(domain.HalResource))
	}
}
