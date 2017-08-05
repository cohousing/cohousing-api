package api

import (
	"github.com/cohousing/cohousing-tenant-api/db"
	"github.com/cohousing/location"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateRouter() {
	router := gin.Default()
	router.Use(ContextResolver())
	router.Use(location.Default())
	router.Use(Authenticate())

	apiV1 := router.Group("api/v1")

	apiV1.POST("/login", Login)

	CreateHomeRoutes(apiV1)
	CreateApartmentRoutes(apiV1, dbFactory)
	CreateResidentRoutes(apiV1, dbFactory)
	CreateUserRoutes(apiV1, dbFactory)
	CreateGroupRoutes(apiV1, dbFactory)

	CreateFixtureRoutes(apiV1)

	router.Run(":8080")
}

func dbFactory(c *gin.Context) *gorm.DB {
	t := GetTenantFromContext(c)
	return db.GetTenantDB(t)
}
