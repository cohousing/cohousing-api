package main

import (
	"github.com/cohousing/cohousing-api/api"
	"github.com/cohousing/cohousing-api/config"
	"github.com/cohousing/cohousing-api/db"
)

func main() {
	db.InitDB()

	config.Configure()

	api.CreateRouter()
}
