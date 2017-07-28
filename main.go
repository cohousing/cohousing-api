package main

import (
	"github.com/cohousing/cohousing-tenant-api/api"
	"github.com/cohousing/cohousing-tenant-api/config"
	"github.com/cohousing/cohousing-tenant-api/db"
)

func main() {
	db.InitDB()

	config.Configure()

	api.CreateRouter()
}
