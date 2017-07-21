package api

import (
	"github.com/gin-gonic/gin"
)

func CreateRouter() {
	router := gin.Default()
	router.Use(ContextResolver())

	apiV1 := router.Group("api/v1")

	CreateApartmentRoutes(apiV1)
	CreateResidentRoutes(apiV1)

	router.Run(":8080")
}
