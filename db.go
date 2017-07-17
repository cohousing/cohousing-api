package main

import (
	"fmt"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

var (
	dbCache map[string]*gorm.DB = make(map[string]*gorm.DB)
	confDB  *gorm.DB
)

func getTenantDB(tenant *Tenant) *gorm.DB {
	db := dbCache[tenant.Context]
	fmt.Fprintf(os.Stdout, "Getting DB connection for %s, got %v\n", tenant.Context, db)
	if db == nil {
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.ConfigDB.User, config.ConfigDB.Password, config.ConfigDB.Host, config.ConfigDB.Port, tenant.Context)
		fmt.Fprintf(os.Stdout, "Establishing db connection for %s: %s\n", tenant.Context, connectionString)
		var err error
		db, err = gorm.Open("mysql", connectionString)
		if err != nil {
			panic(err)
		}
		dbCache[tenant.Context] = db
		db.AutoMigrate(&domain.Apartment{}, &domain.Resident{})
	}
	return db
}

func getConfDB() *gorm.DB {
	if confDB == nil {
		var err error
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.ConfigDB.User, config.ConfigDB.Password, config.ConfigDB.Host, config.ConfigDB.Port, config.ConfigDB.Name)
		confDB, err = gorm.Open("mysql", connectionString)
		if err != nil {
			panic(err)
		}
		confDB.AutoMigrate(&Tenant{})
	}
	return confDB
}
