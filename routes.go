package main

import (
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
)

func CreateRouter() {
	router := gin.Default()
	router.Use(ContextResolver())

	api := router.Group("api/v1")

	api.GET("/", homeRoute)

	ConfigureBasicTenantEndpoint(api, "apartments", domain.Apartment{})
	ConfigureBasicTenantEndpoint(api, "residents", domain.Resident{})

	router.Run(":8080")
}
